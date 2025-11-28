<template>
  <div class="space-y-4">
    <!-- File Input Area -->
    <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-6 text-center hover:border-primary transition cursor-pointer"
         @click="openFileDialog"
         @drop.prevent="handleDrop"
         @dragover.prevent>
      <input ref="fileInput" type="file" multiple @change="handleFileSelect" class="hidden" :accept="acceptedFormats" />
      <div class="text-4xl mb-2">📸</div>
      <p class="font-semibold">Drop files here or click to browse</p>
      <p class="text-sm text-gray-600 dark:text-gray-400">Supported: JPG, PNG, PDF (max 10MB each)</p>
    </div>

    <!-- Upload Progress -->
    <div v-if="uploading" class="space-y-2">
      <p class="text-sm font-semibold">Uploading...</p>
      <div class="bg-gray-200 dark:bg-gray-700 rounded-full h-2 overflow-hidden">
        <div class="bg-primary h-full transition-all" :style="{ width: uploadProgress + '%' }"></div>
      </div>
    </div>

    <!-- File List -->
    <div v-if="uploadedFiles.length" class="space-y-2">
      <p class="text-sm font-semibold">Attached Files</p>
      <div v-for="file in uploadedFiles" :key="file.id" class="flex items-center justify-between bg-gray-100 dark:bg-gray-800 p-3 rounded">
        <div class="flex items-center space-x-2">
          <span class="text-xl">📎</span>
          <span class="text-sm">{{ file.filename }}</span>
        </div>
        <button @click="deleteFile(file.id)" class="text-red-500 hover:text-red-700 text-lg">×</button>
      </div>
    </div>

    <!-- Status Messages -->
    <div v-if="error" class="bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-100 p-3 rounded">
      {{ error }}
    </div>
    <div v-if="success" class="bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-100 p-3 rounded">
      {{ success }}
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const props = defineProps({
  entryType: {
    type: String,
    required: true,
  },
  entryId: {
    type: [String, Number],
    required: true,
  },
})

const emit = defineEmits(['uploaded'])

const fileInput = ref(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadedFiles = ref([])
const error = ref('')
const success = ref('')
const acceptedFormats = '.jpg,.jpeg,.png,.pdf'

const openFileDialog = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  uploadFiles(event.target.files)
}

const handleDrop = (event) => {
  uploadFiles(event.dataTransfer.files)
}

const uploadFiles = (files) => {
  for (const file of files) {
    uploadFile(file)
  }
}

const uploadFile = (file) => {
  if (file.size > 10 * 1024 * 1024) {
    error.value = `${file.name} is too large (max 10MB)`
    return
  }

  const formData = new FormData()
  formData.append('file', file)
  formData.append('entry_type', props.entryType)
  formData.append('entry_id', props.entryId)

  uploading.value = true
  error.value = ''
  success.value = ''

  const xhr = new XMLHttpRequest()

  xhr.upload.addEventListener('progress', (event) => {
    if (event.lengthComputable) {
      uploadProgress.value = Math.round((event.loaded / event.total) * 100)
    }
  })

  xhr.addEventListener('load', () => {
    if (xhr.status === 201) {
      const response = JSON.parse(xhr.responseText)
      uploadedFiles.value.push(response.attachment)
      success.value = `${file.name} uploaded successfully`
      emit('uploaded', response.attachment)
    } else {
      const response = JSON.parse(xhr.responseText)
      error.value = response.error || 'Upload failed'
    }
    uploading.value = false
    uploadProgress.value = 0
  })

  xhr.addEventListener('error', () => {
    error.value = 'Upload failed'
    uploading.value = false
  })

  xhr.open('POST', 'http://localhost:3000/api/upload')
  xhr.setRequestHeader('Authorization', authStore.token)
  xhr.send(formData)
}

const deleteFile = async (fileId) => {
  try {
    const response = await fetch(`http://localhost:3000/api/attachments/${fileId}`, {
      method: 'DELETE',
      headers: { 'Authorization': authStore.token },
    })

    if (response.ok) {
      uploadedFiles.value = uploadedFiles.value.filter(f => f.id !== fileId)
      success.value = 'File deleted'
    } else {
      error.value = 'Failed to delete file'
    }
  } catch (err) {
    error.value = 'Delete error'
  }
}
</script>
