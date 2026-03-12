<script lang="ts" setup>
  import HelloWorld from './components/HelloWorld.vue'
  import MainMenuSideBar from './components/MainMenuSideBar.vue'
</script>

<template>
  <MainMenuSideBar/>
  <main>
    <HelloWorld/>
  </main>
</template>

<style>
main {
  display: flex;
}

</style>
