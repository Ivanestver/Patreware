<script lang="ts" setup>
import {reactive, ref} from 'vue'
import { GetDirPathToScan, GetFilePathToScan, StartScan } from "../../wailsjs/go/main/App"
import { EventsOn } from "../../wailsjs/runtime"
import { main } from "../../wailsjs/go/models"

const isVirusStruct = reactive({
    isVirusFile: false,
    isVirusDir: false
})

const inputPath = reactive({
    toFile: "",
    toDir: ""
})

class resultViewModel {
    path: string;
    isVirus: boolean;

    constructor(path: string, isVirus: boolean) {
        this.path = path;
        this.isVirus = isVirus;
    }

    getPath() {
        return this.path;
    }

    getIsVirus() {
        return this.isVirus
    }
}
const resultsViewModel = ref(Array<resultViewModel>())

const percentCompleted = ref(0)
const checkHasRun = ref(false)

async function getFileToCheck() {
    const path = await GetFilePathToScan()
    if (path) {
        inputPath.toFile = path
    }
}

async function getDirectoryToCheck() {
    const path = await GetDirPathToScan()
    if (path) {
        inputPath.toDir = path
    }
}

async function checkFile() {
    checkHasRun.value = true
    resultsViewModel.value.length = 0
    await StartScan(inputPath.toFile)
}

async function checkDir() {
    checkHasRun.value = true
    resultsViewModel.value.length = 0
    await StartScan(inputPath.toDir)
}

EventsOn('scan_progress', (event: main.UIScanEvent) => {
    // percentCompleted.value = event.progress_percent !== undefined ? event.progress_percent : 0
    // isVirusStruct.isVirusFile = event.virus_found !== undefined && event.virus_found
    resultsViewModel.value.push(new resultViewModel(
        event.current_file ?? "",
        event.virus_found ?? false
    ))
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
                    <input :value="inputPath.toFile" type="text" name="file-path" id="exact-scanning-file-field" class="exact-scanning-file-field">
                    <button type="button" v-on:click="getFileToCheck">Выбрать файл</button>
                </div>
                <button type="button" v-on:click="checkFile">Запустить сканирование</button>
                <div v-show="checkHasRun">
                    <span>{{ percentCompleted }}%</span>
                </div>
            </div>
            <div class="divider"></div>
            <div class="exact-scanning-file-form">
                <div class="exact-scanning-file-div">
                    <label for="exact-scanning-file-field" class="exact-scanning-file-label">Выбрать директорию</label>
                    <input :value="inputPath.toDir" type="text" name="file-path" id="exact-scanning-file-field" class="exact-scanning-file-field">
                    <button type="button" v-on:click="getDirectoryToCheck">Выбрать директорию</button>
                </div>
                <button type="button" v-on:click="checkDir">Запустить сканирование</button>
                <div v-show="checkHasRun">
                    <span>{{ percentCompleted }}%</span>
                </div>
            </div>
            <table>
                <thead>
                    <th>Файл</th>
                    <th>Результат проверки</th>
                </thead>
                <tbody>
                    <tr v-for="res in resultsViewModel">
                        <td>{{ res.getPath() }}</td>
                        <td v-if="res.getIsVirus()">Заражён</td>
                        <td v-else>Не заражён</td>
                    </tr>
                </tbody>
            </table>
        </section>
    </div>
</template>

<style scoped>
.scanning {
    width: 100%;
    height: 100%;
}

.divider {
    height: 0px;
    width: 100%;
    border: 1px solid black;
}
</style>
