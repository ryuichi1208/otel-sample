<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Vue TEST</title>
<!-- Vue.js を読み込む -->
<script src="https://cdn.jsdelivr.net/npm/vue"></script>
</head>
<body>

<div id="app-1">{{ message }}</div>  <!-- {{ message }} が Vueデータに置換される -->

<script>
var app1 = new Vue({
  el: '#app-1',                        /* #app-1 要素に対して Vue を適用する */
  data: { message: 'Hello world!' }    /* message という名前のデータを定義する */
})
</script>

</body>
</html>

