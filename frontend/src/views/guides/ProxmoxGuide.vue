<template>
  <div class="guide-container">
    <h1>Proxmox Integration Guide</h1>

    <div class="guide-steps">
      <section class="step">
        <h2>1. Create a Receiver</h2>
        <p>First, create a new receiver in MNLRNotify:</p>
        <ul>
          <li>Go to Receivers section</li>
          <li>Click "New Receiver"</li>
          <li>Select "Gotify Compatible"</li>
          <li>Note down the token and endpoint URL</li>
        </ul>
      </section>

      <section class="step">
        <h2>2. Configure Proxmox</h2>
        <h3>Method A: Using the GUI (Recommended)</h3>
        <ol>
          <li>Navigate to Datacenter > Notifications</li>
          <li>Click "Add" to create a new notification target</li>
          <li>Select "Gotify" from the target type dropdown</li>
          <li>Enter your MNLRNotify token</li>
          <li>
            For the URL, enter your MNLRNotify instance address (e.g.,
            https://example.com) -
            <strong>Important: Do not include a trailing slash</strong>
          </li>
          <li>Click "Test" to verify the connection works</li>
        </ol>

        <h3>Method B: Using the Command Line</h3>
        <div class="code-block">
          <pre><code>
# Edit notification configuration
nano /etc/pve/notify.cfg

# Add the following section:
gotify:
  token your-mnlr-token
  url https://your-mnlr-instance/message
          </code></pre>
        </div>
      </section>

      <section class="step">
        <h2>3. Test Configuration</h2>
        <p>Test the notification setup using either:</p>
        <ul>
          <li>The "Test" button in the GUI</li>
          <li>Or via Proxmox CLI:</li>
        </ul>
        <div class="code-block">
          <pre><code>
pvesh create /nodes/localhost/test_notification --type gotify
          </code></pre>
        </div>
      </section>

      <BAlert show variant="info" class="mt-4">
        <h4>Note</h4>
        <p>
          Ensure your MNLRNotify instance is reachable from your Proxmox server
          and that you have configured appropriate routing in your channel
          settings. If you receive confusing errors, double-check that your URL
          doesn't have a trailing slash.
        </p>
      </BAlert>
    </div>
  </div>
</template>

<style scoped>
.guide-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.step {
  margin-bottom: 2rem;
  padding: 1rem;
  background: var(--bs-light-bg-subtle);
  border-radius: 8px;
}

.code-block {
  background: var(--bs-dark-bg-subtle);
  padding: 1rem;
  border-radius: 4px;
  margin: 1rem 0;
}

.code-block pre {
  margin: 0;
}

.code-block code {
  color: var(--bs-emphasis-color);
}
</style>
