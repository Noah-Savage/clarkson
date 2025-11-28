export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">Clarkson</h1>
        <p className="text-lg text-muted-foreground mb-8">Vehicle Expense Tracker with Maintenance Reminders</p>

        <div className="bg-card border rounded-lg p-8 max-w-2xl">
          <h2 className="text-2xl font-semibold mb-4">Getting Started</h2>
          <p className="text-muted-foreground mb-6">
            Clarkson has been exported as a standalone full-stack application. This Next.js project was used only for
            preview purposes.
          </p>

          <div className="space-y-4 text-left">
            <div>
              <h3 className="font-semibold mb-2">Full-Stack Application</h3>
              <p className="text-sm text-muted-foreground">
                Go backend with Gin/GORM/SQLite, Vue.js frontend, Docker deployment ready
              </p>
            </div>

            <div>
              <h3 className="font-semibold mb-2">GitHub & Unraid Ready</h3>
              <p className="text-sm text-muted-foreground">
                Complete source code structure with docker-compose.yml for Unraid deployment
              </p>
            </div>

            <div>
              <h3 className="font-semibold mb-2">Key Features</h3>
              <ul className="text-sm text-muted-foreground list-disc list-inside">
                <li>Multi-vehicle & multi-user management</li>
                <li>Fuel & expense tracking with photos</li>
                <li>Maintenance reminders with auto-detection</li>
                <li>Dark mode & PWA support</li>
                <li>Reports & CSV/JSON export</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
