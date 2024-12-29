<script setup lang="ts">
import TelegramChatIdHelper from "@/components/guides/telegramChatIdHelper.vue";
</script>

<template>
  <div class="guide-container">
    <h1>Telegram Bot Token Guide</h1>

    <div class="guide-steps">
      <section class="step">
        <h2>1. Start a Chat with BotFather</h2>
        <p>Open Telegram and start by creating a new bot:</p>
        <ul>
          <li>Search for "@BotFather" in Telegram</li>
          <li>Start a chat with BotFather by clicking "Start"</li>
        </ul>
      </section>

      <section class="step">
        <h2>2. Create a New Bot</h2>
        <ol>
          <li>Send the command /newbot to BotFather</li>
          <li>Choose a name for your bot (e.g., "My Notify Bot")</li>
          <li>Choose a username for your bot (must end in 'bot')</li>
          <li>BotFather will provide you with a token - save this securely</li>
        </ol>
      </section>

      <section class="step">
        <h2>3. Configure Bot Privacy</h2>
        <p>Optional but recommended steps:</p>
        <ol>
          <li>Send /setprivacy to BotFather</li>
          <li>Select your bot</li>
          <li>
            Choose "Enable" to ensure your bot only receives commands it's meant
            to handle
          </li>
        </ol>
      </section>

      <BAlert :model-value="true" variant="info" class="mt-4">
        <h4>Important Security Notes</h4>
        <ul>
          <li>Keep your bot token secure - it's like a password</li>
          <li>Never share your bot token publicly</li>
          <li>
            If your token is compromised, use /revoke with BotFather to get a
            new one
          </li>
        </ul>
      </BAlert>

      <TelegramChatIdHelper />
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
</style>
