<template>
  <div class="editor-page-container">
    <MonacoTreeEditor
      :font-size="14"
      :files="files"
      :sider-min-width="240"
      filelist-title="Mutations"
      @reload="handleReload"
      @new-file="handleNewFile"
      @new-folder="handleNewFolder"
      @save-file="handleSaveFile"
      @delete-file="handleDeleteFile"
      @delete-folder="handleDeleteFolder"
      @rename-file="handleRename"
      @rename-folder="handleRename"
    ></MonacoTreeEditor>
  </div>
</template>

<script setup lang="ts">
import * as monaco from "monaco-editor";
// @ts-expect-error MonacoEnvironment is not defined
import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
// @ts-expect-error MonacoEnvironment is not defined
import cssWorker from "monaco-editor/esm/vs/language/css/css.worker?worker";
// @ts-expect-error MonacoEnvironment is not defined
import htmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker";
// @ts-expect-error MonacoEnvironment is not defined
import jsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker";
// @ts-expect-error MonacoEnvironment is not defined
import tsWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker";

self.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === "json") {
      return new jsonWorker();
    }
    if (label === "css" || label === "scss" || label === "less") {
      return new cssWorker();
    }
    if (label === "html" || label === "handlebars" || label === "razor") {
      return new htmlWorker();
    }
    if (label === "typescript" || label === "javascript") {
      return new tsWorker();
    }
    return new editorWorker();
  },
};

// Ignore the following lines, they are for monaco-tree-editor

import { usePb } from "@/composeables/usePb";
import type { ChannelMutationsRecord } from "@/types/pocketbase-types";
import { ChannelMutationsExecutorOptions } from "@/types/pocketbase-types";
import { useToastController } from "bootstrap-vue-next";

import {
  Editor as MonacoTreeEditor,
  useMonaco,
  type Files,
} from "monaco-tree-editor";
import { onMounted, ref } from "vue";

const pb = usePb();
const toast = useToastController();
const files = ref<Files>({});

// Example files
const exampleFiles: Files = {
  "VIRT:\\examples\\JavaScript Examples\\uppercase.js": {
    isFile: true,
    content: "// Convert message to uppercase\nmessage.toUpperCase();",
    readonly: true,
  },
  "VIRT:\\examples\\JavaScript Examples\\functions.js": {
    isFile: true,
    content: `
const limitMessageLines = () => {
var temp = message.split("\\n");
var maxLines = 25;
const newMessageArray = temp.slice(0,maxLines)
return newMessageArray.join("\\n")
}

limitMessageLines()
`,
    readonly: true,
  },
  "VIRT:\\examples\\JavaScript Examples\\json-manipulation.js": {
    isFile: true,
    content:
      "// JSON manipulation example\nlet data = JSON.parse(message);\ndata.processed = true;\nJSON.stringify(data);",
    readonly: true,
  },
  "VIRT:\\examples\\Shell Examples\\basic.cmd": {
    isFile: true,
    content: '# Basic echo example\necho "PREFIX: {{message}}"',
    readonly: true,
  },
  "VIRT:\\examples\\PowerShell Examples\\json-processing.cmd": {
    isFile: true,
    content:
      "# JSON processing example\npowershell -noprofile -command \"&{\n    $msg = '{{message}}' | ConvertFrom-Json\n    $msg.level = 'INFO'\n    $msg | ConvertTo-Json\n}\"",
    readonly: true,
  },
  "VIRT:\\examples\\JavaScript Examples\\kvStore.js": {
    isFile: true,
    content: `// KV Store Example
// Store data
const userPrefs = { theme: "dark", notifications: true };
kvStore.Set("userPreferences", JSON.stringify(userPrefs));

// Retrieve data
const savedPrefs = kvStore.Get("userPreferences");
let prefs = savedPrefs ? JSON.parse(savedPrefs) : {};

// Modify data
prefs.lastMessage = message;

// Store updated data
kvStore.Set("userPreferences", JSON.stringify(prefs));

// Return the processed message
\`Processed by \${prefs.theme} theme: \${message}\`;`,
    readonly: true,
  },
};

// Convert mutation to filename
const getMutationFilename = (mutation: ChannelMutationsRecord) => {
  const ext =
    mutation.executor === ChannelMutationsExecutorOptions.goja ? ".js" : ".cmd";
  return `${mutation.name}${ext}`;
};

const getMutationNameFromPath = (path: string) => {
  return (
    path
      .split("\\")
      .pop()
      ?.replace(/\.(js|cmd)$/, "") || ""
  );
};

// Convert mutations to file structure
const convertMutationsToFiles = (mutations: ChannelMutationsRecord[]) => {
  const newFiles: Files = { ...exampleFiles };
  mutations.forEach((mutation) => {
    newFiles[`VIRT:\\mutations\\${getMutationFilename(mutation)}`] = {
      isFile: true,
      content: mutation.script,
    };
  });
  return newFiles;
};

// Load mutations from PocketBase
const loadMutations = async () => {
  try {
    const response = await pb.collection("channel_mutations").getFullList();
    files.value = convertMutationsToFiles(response);
  } catch (error) {
    console.error("Failed to load mutations", error);
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to load mutations",
        variant: "danger",
      },
    });
  }
};

const handleReload = async (
  resolve: () => void,
  reject: (msg?: string) => void
) => {
  try {
    await loadMutations();
    resolve();
  } catch (error) {
    reject("Failed to reload mutations");
  }
};

const handleSaveFile = async (
  path: string,
  content: string,
  resolve: () => void,
  reject: (msg?: string) => void
) => {
  try {
    if (path.startsWith("VIRT:\\examples\\")) {
      reject("Example files cannot be modified");
      return;
    }
    const mutationName = getMutationNameFromPath(path);
    const records = await pb.collection("channel_mutations").getFullList({
      filter: `name = "${mutationName}"`,
    });
    if (records.length === 0) throw new Error("Mutation not found");

    await pb.collection("channel_mutations").update(records[0].id, {
      script: content,
    });
    resolve();
  } catch (error) {
    reject("Failed to save mutation");
  }
};

const handleNewFile = async (
  path: string,
  resolve: Function,
  reject: Function
) => {
  try {
    const filename = path.split("\\").pop() || "";
    const isCmd = filename.endsWith(".cmd");

    const mutation = await pb.collection("channel_mutations").create({
      name: filename.replace(/\.(js|cmd)$/, ""),
      script: filename.endsWith(".cmd")
        ? 'echo "PREFIX: {{message}}"'
        : "message",
      executor: isCmd
        ? ChannelMutationsExecutorOptions.cmd
        : ChannelMutationsExecutorOptions.goja,
    });

    await loadMutations();
    resolve();
  } catch (error) {
    reject("Failed to create mutation");
  }
};

const handleDeleteFile = async (
  path: string,
  resolve: () => void,
  reject: (msg?: string) => void
) => {
  try {
    if (path.startsWith("VIRT:\\examples\\")) {
      reject("Example files cannot be deleted");
      return;
    }
    const mutationName = getMutationNameFromPath(path);
    const records = await pb.collection("channel_mutations").getFullList({
      filter: `name = "${mutationName}"`,
    });
    if (records.length === 0) throw new Error("Mutation not found");

    await pb.collection("channel_mutations").delete(records[0].id);
    await loadMutations();
    resolve();
  } catch (error) {
    reject("Failed to delete mutation");
  }
};

const handleRename = async (
  path: string,
  newPath: string,
  resolve: () => void,
  reject: (msg?: string) => void
) => {
  try {
    if (path.startsWith("VIRT:\\examples\\")) {
      reject("Example files cannot be renamed");
      return;
    }
    const oldName = getMutationNameFromPath(path);
    const newName = getMutationNameFromPath(newPath);

    const records = await pb.collection("channel_mutations").getFullList({
      filter: `name = "${oldName}"`,
    });
    if (records.length === 0) throw new Error("Mutation not found");

    await pb.collection("channel_mutations").update(records[0].id, {
      name: newName,
    });

    await loadMutations();
    resolve();
  } catch (error) {
    reject("Failed to rename mutation");
  }
};

// Disable folder operations
const handleNewFolder = (_: string, __: Function, reject: Function) => {
  reject("Folders are not supported");
};

const handleDeleteFolder = (
  _: string,
  __: () => void,
  reject: (msg?: string) => void
) => {
  reject("Folders are not supported");
};

onMounted(() => {
  //loadMutations();
});

// ================ init monaco-tree-editor ================
monaco.languages.typescript.javascriptDefaults.addExtraLib(
  `
declare let message: string;

declare const kvStore: {
  /**
   * Get a value from the key-value store
   * @param key The key to retrieve
   * @returns The stored value or null if not found
   */
  Get(key: string): string | null;
  
  /**
   * Store a value in the key-value store
   * @param key The key to store under
   * @param value The value to store must be a json string (Use JSON.stringify)
   * @example
   * kvStore.Set("key", JSON.stringify({ foo: "bar" }));
   */
  Set(key: string, value: string): void;
  
  /**
   * Delete a value from the key-value store
   * @param key The key to delete
   */
  Delete(key: string): void;
};
`,
  "global.d.ts"
);
monaco.languages.typescript.javascriptDefaults.setCompilerOptions({
  allowNonTsExtensions: true,
  target: monaco.languages.typescript.ScriptTarget.ES5,
  lib: ["es5"],
});

const monacoStore = useMonaco(monaco);
</script>

<style>
@import "monaco-tree-editor/index.css";

@media (max-width: 768px) {
  .editor-page-container {
    height: calc(100vh - 76px);
  }
}

@media (min-width: 768px) {
  .editor-page-container {
    height: calc(100vh - 56px);
  }
}

.editor-page-container {
  display: flex;
  flex-direction: column;
}
</style>
