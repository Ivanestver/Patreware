<script lang="ts" setup>
import {reactive, ref} from 'vue'
import { GetPathToScan, StartScan } from "../../wailsjs/go/main/App"
import { EventsOn } from "../../wailsjs/runtime"
import { main } from "../../wailsjs/go/models"

const inputPath = ref("")
const isVirus = ref(false)

async function getFileToCheck() {
    const path = await GetPathToScan()
    if (path) {
        inputPath.value = path
    }
}

async function check() {
    await StartScan(inputPath.value)
}

EventsOn('scan_progress', (event: main.UIScanEvent) => {
    isVirus.value = event.virus_found === undefined ? false : event.virus_found
})

</script>

<template>
    <div class="scanning">
        <section class="exact-scanning">
            <p class="exact-scanning-desc">
                Проведение сканирования определённого файла или директории
            </p>
            <div class="exact-scanning-file-form">
                <div class="exact-scanning-file-div">
                    <label for="exact-scanning-file-field" class="exact-scanning-file-label">Выбрать файл</label>
                    <input :value="inputPath" type="text" name="file-path" id="exact-scanning-file-field" class="exact-scanning-file-field">
                    <button type="button" v-on:click="getFileToCheck">Выбрать файл</button>
                </div>
                <button type="button" v-on:click="check">Запустить сканирование</button>
                <span v-if="isVirus">Заражён</span>
                <span v-else>Не заражён</span>
            </div>
        </section>
    </div>
</template>

<style scoped>
</style>
