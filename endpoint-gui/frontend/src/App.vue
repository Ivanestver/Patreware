<script lang="ts" setup>
  import { ref } from "vue"
  import Scanning from './components/Scanning.vue'
  import MainMenuSideBar from './components/MainMenuSideBar.vue'

  enum Tab {ScanningTab, InfoTab};

  const currentTab = ref(Tab.ScanningTab)

</script>

<template>
  <MainMenuSideBar/>
  <main>
    <img class="main-logo" src="./assets/images/main_logo.jpeg" alt="main-logo" />
    <Scanning v-show="currentTab === Tab.ScanningTab" />
  </main>
</template>

<style>
main {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  position: relative;
}

.main-logo {
  position: absolute;
  z-index: -1;
  width: 100%;
  height: 100%;
  opacity: 0.3;
}

</style>
