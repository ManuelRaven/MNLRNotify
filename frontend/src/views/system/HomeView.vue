<script setup lang="ts"></script>

<template>
  <main>
    <div class="app-tutorial">
      <h1>Welcome to MNLRNotify</h1>
      <p class="slogan">
        A focused message routing system with basic Gotify/ntfy.sh compatibility
      </p>

      <section class="description">
        <h2>About MNLRNotify</h2>
        <p>
          MNLRNotify implements the core message handling endpoints of Gotify
          and ntfy.sh, focusing solely on message reception and routing. It
          provides basic API compatibility for sending messages, but does not
          support advanced features like user management, applications, or
          websocket connections found in the original services.
        </p>
      </section>

      <section class="compatibility-notice">
        <h2>Compatibility Notice</h2>
        <p>
          This application only implements the basic message POST endpoints:
        </p>
        <ul>
          <li>Basic message reception (POST /message)</li>
          <li>Simple authentication</li>
          <li>Basic priority handling</li>
        </ul>
        <p class="warning">
          Advanced features from Gotify/ntfy.sh such as user management,
          websockets, applications, or subscription management are not
          supported.
        </p>
      </section>

      <section class="message-flow">
        <h2>Message Flow</h2>
        <div class="flow-diagram">
          <div class="flow-step">
            <h3>1. Receiver</h3>
            <p>Compatible with Gotify/Ntfy API formats</p>
            <small>Use existing clients without modifications</small>
          </div>
          <div class="flow-arrow">→</div>
          <div class="flow-step">
            <h3>2. Channel</h3>
            <p>Advanced message routing</p>
            <small>Controls message flow and delivery rules</small>
          </div>
          <div class="flow-arrow">→</div>
          <div class="flow-step">
            <h3>3. Sender</h3>
            <p>Multiple delivery options</p>
            <small>Forward to any Shoutrrr-supported service</small>
          </div>
        </div>
      </section>

      <section class="endpoints">
        <h2>Limited API Compatibility</h2>
        <div class="endpoint-example">
          <h3>Basic Message Endpoint</h3>
          <p>Only the core message posting functionality is supported:</p>
          <code class="endpoint-code">
            POST /message<br />
            {<br />
            "message": "Hello World",<br />
            "priority": 5,<br />
            "title": "Test Message"<br />
            }
          </code>
        </div>
      </section>

      <section class="quick-start">
        <h2>Quick Start Guide</h2>
        <ol>
          <li>Create a Channel to handle message routing</li>
          <li>Set up Receivers to accept Gotify/Ntfy messages</li>
          <li>Add Senders with Shoutrrr URLs</li>
          <li>Connect Receivers to Channels</li>
          <li>Connect Channels to Senders</li>
        </ol>
      </section>
    </div>
  </main>
</template>

<style scoped>
.app-tutorial {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.flow-diagram {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin: 2rem 0;
  padding: 1rem;
  background: var(--bs-dark-bg-subtle);
  border-radius: 8px;
}

.flow-step {
  padding: 1rem;
  background: var(--bs-tertiary-bg);
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  flex: 1;
}

.flow-arrow {
  color: var(--bs-secondary);
  font-weight: bold;
  font-size: 1.5rem;
}

.endpoint-example {
  background: var(--bs-dark-bg-subtle);
  padding: 1rem;
  border-radius: 6px;
  margin: 1rem 0;
}

.endpoint-code {
  display: block;
  background: var(--bs-tertiary-bg);
  color: var(--bs-emphasis-color);
  padding: 1rem;
  border-radius: 4px;
  font-family: monospace;
  margin: 1rem 0;
}

section {
  margin-bottom: 2rem;
}

h1 {
  color: var(--bs-emphasis-color);
  margin-bottom: 2rem;
}

h2 {
  color: var(--bs-primary);
  margin-bottom: 1rem;
}

h3 {
  color: var(--bs-emphasis-color);
}

p,
small {
  color: var(--bs-secondary-color);
}

code {
  background: var(--bs-tertiary-bg);
  color: var(--bs-emphasis-color);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
}

.quick-start ol {
  padding-left: 1.5rem;
  color: var(--bs-body-color);
}

.quick-start li {
  margin: 0.5rem 0;
}

.slogan {
  font-size: 1.25rem;
  color: var(--bs-secondary-color);
  text-align: center;
  margin-bottom: 2.5rem;
  font-style: italic;
}

.compatibility-notice {
  background: var(--bs-warning-bg-subtle);
  border-radius: 8px;
  padding: 1rem;
  margin: 2rem 0;
}

.compatibility-notice ul {
  padding-left: 1.5rem;
  margin: 1rem 0;
  color: var(--bs-body-color);
}

.warning {
  color: var(--bs-warning-text);
  font-weight: 500;
  margin-top: 1rem;
}
</style>
