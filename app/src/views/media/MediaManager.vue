<script setup lang="ts">
import type { MediaEntry, MediaRoot } from '@/api/media'
import {
  DeleteOutlined,
  FolderAddOutlined,
  FolderOpenOutlined,
  ReloadOutlined,
  UploadOutlined,
} from '@antdv-next/icons'
import { useDropZone } from '@vueuse/core'
import { App, Button, Card, Empty, Input, Space, Spin, Tag } from 'antdv-next'
import dayjs from 'dayjs'
import mediaApi from '@/api/media'
import { bytesToSize } from '@/lib/helper'

const { message, modal } = App.useApp()

const loading = ref(false)
const uploading = ref(false)
const roots = ref<MediaRoot[]>([])
const selectedRoot = ref('')
const currentPath = ref('')
const entries = ref<MediaEntry[]>([])
const selected = ref<string[]>([])
const previewMap = ref<Record<string, string>>({})
const fileInputRef = useTemplateRef<HTMLInputElement>('fileInputRef')
const dropZoneRef = useTemplateRef<HTMLElement>('dropZoneRef')

const breadcrumbs = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  const items = [{ title: selectedRoot.value || $gettext('Root'), path: '' }]
  let acc = ''
  for (const part of parts) {
    acc = acc ? `${acc}/${part}` : part
    items.push({ title: part, path: acc })
  }
  return items
})

const allSelected = computed({
  get: () => {
    const selectable = entries.value.map(item => item.path)
    return selectable.length > 0 && selectable.every(path => selected.value.includes(path))
  },
  set: (checked: boolean) => {
    selected.value = checked ? entries.value.map(item => item.path) : []
  },
})

const { isOverDropZone } = useDropZone(dropZoneRef, {
  onDrop: (files) => {
    if (!files?.length || !selectedRoot.value)
      return
    void handleUpload(Array.from(files))
  },
})

async function loadRoots() {
  loading.value = true
  try {
    const res = await mediaApi.listRoots()
    roots.value = Array.isArray(res?.data) ? res.data : []
    if (!selectedRoot.value && roots.value.length > 0)
      selectedRoot.value = roots.value[0].name
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to load media folders'))
  }
  finally {
    loading.value = false
  }
}

async function loadEntries() {
  if (!selectedRoot.value)
    return
  loading.value = true
  selected.value = []
  try {
    const res = await mediaApi.list(selectedRoot.value, currentPath.value)
    const list = Array.isArray(res?.entries) ? res.entries : []
    entries.value = [...list].sort((a, b) => {
      if (a.is_dir !== b.is_dir)
        return a.is_dir ? -1 : 1
      return a.name.localeCompare(b.name)
    })
    await loadPreviews()
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to load folder contents'))
    entries.value = []
  }
  finally {
    loading.value = false
  }
}

async function loadPreviews() {
  Object.values(previewMap.value).forEach((url) => {
    if (url.startsWith('blob:'))
      URL.revokeObjectURL(url)
  })
  previewMap.value = {}

  const images = entries.value.filter(item => item.is_image).slice(0, 80)
  await Promise.all(images.map(async (item) => {
    try {
      const blob = await mediaApi.fetchPreviewBlob(selectedRoot.value, item.path)
      if (blob instanceof Blob)
        previewMap.value[item.path] = URL.createObjectURL(blob)
    }
    catch {
      return
    }
  }))
}

function openRoot(name: string) {
  selectedRoot.value = name
  currentPath.value = ''
  void loadEntries()
}

function openEntry(entry: MediaEntry) {
  if (!entry.is_dir)
    return
  currentPath.value = entry.path
  void loadEntries()
}

function goToPath(path: string) {
  currentPath.value = path
  void loadEntries()
}

function isChecked(path: string) {
  return selected.value.includes(path)
}

function setChecked(path: string, checked: boolean) {
  if (checked) {
    if (!selected.value.includes(path))
      selected.value = [...selected.value, path]
    return
  }
  selected.value = selected.value.filter(item => item !== path)
}

function pickFiles() {
  fileInputRef.value?.click()
}

function onFileInputChange(event: Event) {
  const input = event.target as HTMLInputElement
  const files = input.files ? Array.from(input.files) : []
  input.value = ''
  if (files.length)
    void handleUpload(files)
}

async function handleUpload(files: File[]) {
  if (!selectedRoot.value || !files.length)
    return
  uploading.value = true
  try {
    await mediaApi.upload(selectedRoot.value, currentPath.value, files)
    message.success($gettext('Upload complete'))
    await loadEntries()
    await loadRoots()
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Upload failed'))
  }
  finally {
    uploading.value = false
  }
}

function createFolder() {
  let name = ''
  modal.confirm({
    title: $gettext('Add Folder'),
    content: () => h(Input, {
      placeholder: $gettext('Folder name'),
      'onUpdate:value': (value: string) => {
        name = value
      },
    }),
    async onOk() {
      if (!name.trim()) {
        message.warning($gettext('Folder name is required'))
        return Promise.reject(new Error('empty'))
      }
      await mediaApi.mkdir(selectedRoot.value, currentPath.value, name.trim())
      message.success($gettext('Folder created'))
      await loadEntries()
      await loadRoots()
    },
  })
}

function confirmDelete(paths: string[]) {
  if (!paths.length)
    return
  modal.confirm({
    title: $gettext('Delete selected items'),
    content: $gettext('This will permanently delete the selected files or folders.'),
    okType: 'danger',
    okText: $gettext('Delete'),
    async onOk() {
      await mediaApi.delete(selectedRoot.value, paths)
      message.success($gettext('Deleted'))
      await loadEntries()
      await loadRoots()
    },
  })
}

onMounted(async () => {
  await loadRoots()
  if (selectedRoot.value)
    await loadEntries()
})

onUnmounted(() => {
  Object.values(previewMap.value).forEach((url) => {
    if (url.startsWith('blob:'))
      URL.revokeObjectURL(url)
  })
})
</script>

<template>
  <div class="media-manager">
    <Card
      :title="$gettext('Media Manager')"
      variant="borderless"
      class="mb-4"
    >
      <template #extra>
        <Space wrap>
          <Button :loading="loading" @click="loadEntries">
            <template #icon>
              <ReloadOutlined />
            </template>
            {{ $gettext('Refresh') }}
          </Button>
          <Button :disabled="!selectedRoot" @click="createFolder">
            <template #icon>
              <FolderAddOutlined />
            </template>
            {{ $gettext('Add Folder') }}
          </Button>
          <Button
            type="primary"
            :loading="uploading"
            :disabled="!selectedRoot"
            @click="pickFiles"
          >
            <template #icon>
              <UploadOutlined />
            </template>
            {{ $gettext('Upload') }}
          </Button>
          <Button
            danger
            :disabled="!selected.length"
            @click="confirmDelete(selected)"
          >
            <template #icon>
              <DeleteOutlined />
            </template>
            {{ $gettext('Delete Selected') }}
          </Button>
        </Space>
      </template>

      <input
        ref="fileInputRef"
        type="file"
        class="hidden"
        multiple
        accept="image/*,.webp,.png,.jpg,.jpeg,.gif,.avif"
        @change="onFileInputChange"
      >

      <div class="mb-4 flex flex-wrap gap-2">
        <Tag
          v-for="root in roots"
          :key="root.name"
          class="cursor-pointer px-3 py-1"
          :color="selectedRoot === root.name ? 'blue' : undefined"
          @click="openRoot(root.name)"
        >
          {{ root.name }}
          <span class="ml-2 opacity-70">{{ root.count }} · {{ bytesToSize(root.size) }}</span>
        </Tag>
        <Empty
          v-if="!roots.length && !loading"
          class="w-full"
          :description="$gettext('No media folders found under /var/www. Mount or create folders there to manage images.')"
        />
      </div>

      <div v-if="selectedRoot" class="mb-2 flex flex-wrap items-center justify-between gap-3">
        <div class="flex flex-wrap items-center gap-2 text-sm">
          <a
            v-for="item in breadcrumbs"
            :key="item.path || 'root'"
            class="breadcrumb-link"
            @click.prevent="goToPath(item.path)"
          >
            {{ item.title }}
            <span class="opacity-40 mx-1">/</span>
          </a>
        </div>
        <ACheckbox v-model:checked="allSelected" :disabled="!entries.length">
          {{ $gettext('Select all') }}
        </ACheckbox>
      </div>
    </Card>

    <div
      v-if="selectedRoot"
      ref="dropZoneRef"
      class="media-dropzone"
      :class="{ 'is-over': isOverDropZone }"
    >
      <Spin :spinning="loading || uploading">
        <div
          v-if="!entries.length"
          class="empty-upload"
        >
          <UploadOutlined class="text-4xl mb-3 opacity-50" />
          <div class="mb-1 font-medium">
            {{ $gettext('Drag and drop images here') }}
          </div>
          <div class="text-sm opacity-70">
            {{ $gettext('Or use Upload to add trading card and comic book images.') }}
          </div>
        </div>

        <div v-else class="media-grid">
          <div
            v-for="entry in entries"
            :key="entry.path"
            class="media-card"
            @click="openEntry(entry)"
          >
            <div class="media-card-check" @click.stop>
              <ACheckbox
                :checked="isChecked(entry.path)"
                @update:checked="(checked: boolean) => setChecked(entry.path, checked)"
              />
            </div>

            <div class="media-thumb">
              <FolderOpenOutlined v-if="entry.is_dir" class="text-4xl" />
              <img
                v-else-if="entry.is_image && previewMap[entry.path]"
                :src="previewMap[entry.path]"
                :alt="entry.name"
              >
              <div v-else class="text-xs opacity-60">
                {{ entry.extension || $gettext('File') }}
              </div>
            </div>

            <div class="media-meta">
              <div class="media-name" :title="entry.name">
                {{ entry.name }}
              </div>
              <div class="media-sub">
                <template v-if="entry.is_dir">
                  {{ $gettext('Folder') }}
                </template>
                <template v-else>
                  {{ bytesToSize(entry.size) }} · {{ dayjs.unix(entry.mod_time).format('YYYY-MM-DD HH:mm') }}
                </template>
              </div>
            </div>
          </div>
        </div>
      </Spin>
    </div>
  </div>
</template>

<style scoped>
.hidden {
  display: none;
}

.breadcrumb-link {
  color: var(--ant-color-text-secondary);
  cursor: pointer;
}

.breadcrumb-link:hover {
  color: var(--ant-color-primary);
}

.media-dropzone {
  min-height: 360px;
  border: 1px dashed var(--ant-color-border);
  border-radius: 8px;
  padding: 16px;
  background: var(--ant-color-bg-container);
  transition: border-color 0.2s ease, background 0.2s ease;
}

.media-dropzone.is-over {
  border-color: var(--ant-color-primary);
  background: var(--ant-color-primary-bg);
}

.empty-upload {
  min-height: 320px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: var(--ant-color-text-secondary);
}

.media-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
}

.media-card {
  position: relative;
  border: 1px solid var(--ant-color-border-secondary);
  border-radius: 8px;
  background: var(--ant-color-bg-container);
  padding: 10px;
  cursor: pointer;
}

.media-card:hover {
  border-color: var(--ant-color-primary);
}

.media-card-check {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 1;
}

.media-thumb {
  height: 120px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: var(--ant-color-fill-tertiary);
  margin-bottom: 8px;
  color: var(--ant-color-text-secondary);
}

.media-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.media-name {
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--ant-color-text);
}

.media-sub {
  font-size: 12px;
  color: var(--ant-color-text-secondary);
  margin-top: 2px;
}
</style>
