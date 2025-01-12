<template>
  <div>
    <BAccordion class="mb-3">
      <BAccordionItem title="Mutation Documentation">
        <h5>Channel Mutations</h5>
        <p>
          Channel mutations allow you to transform messages before they are sent
          to channels. Each mutation consists of a script that processes the
          incoming message.
        </p>

        <h6>Available Executors</h6>

        <h6 class="mt-4">Goja (JavaScript)</h6>
        <p>
          Execute JavaScript code to transform messages. The message is
          available as the 'message' variable. Return the modified message or
          undefined to keep the original.
        </p>
        <BCard class="mb-3">
          <pre class="mb-0"><code>// Example: Convert message to uppercase
message.toUpperCase();

// Example: Add prefix
"[INFO] " + message;

// Example: JSON manipulation
let data = JSON.parse(message);
data.processed = true;
JSON.stringify(data);</code></pre>
        </BCard>

        <h6>Command Line</h6>
        <p>
          Execute shell commands to process messages. Use
          <code v-html="'{{message}}'"></code> placeholder to insert the
          message. The command's output will become the new message.
        </p>

        <BAlert :model-value="true" variant="warning">
          <ul>
            <li>
              This commands running directly on the system, be careful with the
              commands you run.
            </li>
            <li>
              Only available binaries like PowerShell can run when they are
              natively installed. If this app runs inside a container, they are
              <b>not installed by default.</b>
            </li>
          </ul>
        </BAlert>

        <BTabs content-class="mt-3">
          <BTab title="Bash/Shell" active>
            <BCard class="mb-3">
              <pre class="mb-0"><code># Basic echo example
echo "PREFIX: <code v-html="'{{message}}'"></code>"

# Using sed for text replacement
echo "<code v-html="'{{message}}'"></code>" | sed 's/error/ERROR/g'

# JSON processing with jq
echo '<code v-html="'{{message}}'"></code>' | jq '.level = "INFO"'

# Multiple operations
echo "<code v-html="'{{message}}'"></code>" | tr '[:lower:]' '[:upper:]' | sed 's/^/[LOG] /'</code></pre>
            </BCard>
          </BTab>

          <BTab title="PowerShell">
            <BCard class="mb-3">
              <pre class="mb-0"><code># Basic string output
powershell -noprofile -command "&{Write-Output 'PREFIX: <code v-html="'{{message}}'"></code>'}"

# Converting to uppercase
powershell -noprofile -command "&{$msg = '<code v-html="'{{message}}'"></code>'; $msg.ToUpper()}"

# JSON manipulation
powershell -noprofile -command "&{
    $msg = '<code v-html="'{{message}}'"></code>' | ConvertFrom-Json
    $msg.level = 'INFO'
    $msg | ConvertTo-Json
}"

# Multiple operations
powershell -noprofile -command "&{
    $msg = '<code v-html="'{{message}}'"></code>'
    '[LOG] ' + $msg.ToUpper()
}"</code></pre>
            </BCard>
          </BTab>

          <BTab title="PowerShell 7+">
            <BCard class="mb-3">
              <pre class="mb-0"><code># Basic command (pwsh for PowerShell 7+)
pwsh -noprofile -command "&{Write-Output 'PREFIX: <code v-html="'{{message}}'"></code>'}"

# JSON processing with modern features
pwsh -noprofile -command "&{
    $data = @{
        original = '<code v-html="'{{message}}'"></code>'
        processed = $true
        timestamp = (Get-Date).ToString('o')
    }
    $data | ConvertTo-Json
}"

# Pipeline operations
pwsh -noprofile -command "&{
    '<code v-html="'{{message}}'"></code>' | 
        Where-Object { $_ -match 'error' } |
        ForEach-Object { '[ERROR] $_' } |
        Out-String -NoNewline
}"

# Advanced error handling
pwsh -noprofile -command "&{
    try {
        $msg = '<code v-html="'{{message}}'"></code>' | ConvertFrom-Json
        $msg.priority = 'high'
        $msg | ConvertTo-Json
    } catch {
        '[ERROR] Invalid JSON: <code v-html="'{{message}}'"></code>'
    }
}"</code></pre>
            </BCard>
          </BTab>
        </BTabs>
      </BAccordionItem>
    </BAccordion>

    <TextResponseModal
      v-model="showTextResponseModal"
      :title="'Edit Name'"
      :presetText="currentEditMutation?.name"
      @close="showTextResponseModal = false"
      @save="onNameChange"
    ></TextResponseModal>

    <div class="mb-3">
      <div class="d-flex gap-2 mb-2">
        <BFormInput
          v-model="newMutation.name"
          placeholder="Enter Name"
          @keyup.enter="createMutation"
        />
        <BFormSelect
          v-model="newMutation.executor"
          :options="executorOptions"
        />
        <BButton variant="primary" @click="createMutation"
          >Add Mutation</BButton
        >
      </div>
      <BFormTextarea
        v-model="newMutation.script"
        placeholder="Enter script"
        rows="3"
        class="mb-2"
      />
    </div>

    <b-pagination
      v-model="currentPage"
      :total-rows="rowsCount"
      :per-page="perPage"
    />
    <BTable
      stacked="md"
      striped
      label-stacked
      hover
      :items="mutations"
      :fields="sortFields"
    >
      <template #cell(script)="row">
        <BFormTextarea
          v-model="row.item.script"
          rows="3"
          @update:model-value="onScriptChange(row.item)"
        />
      </template>
      <template #cell(executor)="row">
        <BFormSelect
          v-model="row.item.executor"
          :options="executorOptions"
          @update:model-value="onExecutorChange(row.item)"
        />
      </template>
      <template #cell(actions)="row">
        <BButton
          size="sm"
          variant="danger"
          class="me-1"
          @click="onDelete(row.item, row.index)"
        >
          Delete
        </BButton>

        <BButton
          size="sm"
          variant="primary"
          @click="
            () => {
              currentEditMutation = row.item;
              showTextResponseModal = true;
            }
          "
        >
          Edit Name
        </BButton>
      </template>
    </BTable>
  </div>
</template>

<script setup lang="ts">
import TextResponseModal from "@/components/GenericModals/TextResponseModal.vue";
import { usePb } from "@/composeables/usePb";
import type { ChannelMutationsRequest } from "@/types/custom-types";
import type { ChannelMutationsRecord } from "@/types/pocketbase-types";
import { ChannelMutationsExecutorOptions } from "@/types/pocketbase-types";
import {
  useModalController,
  useToastController,
  type TableFieldRaw,
} from "bootstrap-vue-next";
import { onMounted, ref } from "vue";

const showTextResponseModal = ref(false);
const currentEditMutation = ref<ChannelMutationsRecord | null>(null);

const pb = usePb();
const toast = useToastController();
const { confirm } = useModalController();

const perPage = 10;
const currentPage = ref(1);
const rowsCount = ref(0);

const newMutation = ref<ChannelMutationsRequest>({
  name: "",
  script: "",
  executor: ChannelMutationsExecutorOptions.goja,
});

const mutations = ref<ChannelMutationsRecord[]>([]);
const executorOptions = Object.values(ChannelMutationsExecutorOptions);

const onDelete = async (mutation: ChannelMutationsRecord, index: number) => {
  try {
    const value = await confirm?.({
      props: {
        title: "Are you sure?",
        body: "Do you want to delete this mutation? This action cannot be undone.",
      },
    });

    if (!value) return;
    await pb.collection("channel_mutations").delete(mutation.id);
    toast.show?.({
      props: { title: "Success", body: "Mutation Deleted", variant: "success" },
    });

    mutations.value.splice(index, 1);
  } catch (error) {
    toast.show?.({
      props: { title: "Error", body: (error as any) || "", variant: "danger" },
    });
  }
};

const onNameChange = async (newName: string | null) => {
  showTextResponseModal.value = false;
  if (currentEditMutation.value && newName) {
    try {
      await pb
        .collection("channel_mutations")
        .update(currentEditMutation.value.id, {
          name: newName,
        });

      toast.show?.({
        props: {
          title: "Success",
          body: "Name updated successfully",
          variant: "success",
        },
      });

      await fetchMutations();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update name",
          variant: "danger",
        },
      });
    }
  }
};

const onScriptChange = async (mutation: ChannelMutationsRecord) => {
  try {
    await pb.collection("channel_mutations").update(mutation.id, {
      script: mutation.script,
    });
    toast.show?.({
      props: {
        title: "Success",
        body: "Script updated successfully",
        variant: "success",
      },
    });
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to update script",
        variant: "danger",
      },
    });
  }
};

const onExecutorChange = async (mutation: ChannelMutationsRecord) => {
  try {
    await pb.collection("channel_mutations").update(mutation.id, {
      executor: mutation.executor,
    });
    toast.show?.({
      props: {
        title: "Success",
        body: "Executor updated successfully",
        variant: "success",
      },
    });
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to update executor",
        variant: "danger",
      },
    });
  }
};

const createMutation = async () => {
  await pb.collection("channel_mutations").create(newMutation.value);

  newMutation.value = {
    name: "",
    script: "",
    executor: ChannelMutationsExecutorOptions.goja,
  };

  await fetchMutations();
};

const sortFields: Exclude<TableFieldRaw<ChannelMutationsRecord>, string>[] = [
  { key: "name", label: "Name", sortable: true },
  { key: "executor", label: "Executor" },
  { key: "script", label: "Script" },
  { key: "actions", label: "Actions" },
];

const fetchMutations = async () => {
  try {
    let resp = await pb
      .collection("channel_mutations")
      .getList(currentPage.value, perPage);

    mutations.value = resp.items;
    rowsCount.value = resp.totalItems;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch data",
        variant: "danger",
      },
    });
  }
};

onMounted(() => {
  fetchMutations();
});
</script>

<style scoped>
pre {
  background: var(--bs-dark-bg-subtle);
  padding: 1rem;
  border-radius: 4px;
  margin: 1rem 0;
}
</style>
