<template>
  <div>
    <ul>
<li>当前界面包名:{{currentPackage}}</li>
<li>当前界面活动:{{currentActivity}}</li>
</ul>
<el-select v-model="value" filterable placeholder="Select">
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    />
  </el-select>
  </div>
</template>

<script>
import { ElMessage } from "element-plus";
import { handleError } from "vue";
import { GetCurrentActivity } from "../../wailsjs/go/app/App.js";

import { GetDeviceList, GetDeviceProp, Excute, GetImage,ExcuteSync } from "../../wailsjs/go/app/App.js";
import useGlobalProperties from '../hooks/globalVar.js'


export default {
  data() {
    return {
      currentPackage: "",
      currentActivity:""
    };
  },


  mounted() {
    const globalProperties = useGlobalProperties()
    console.log(globalProperties.$deviceid)

    GetCurrentActivity().then((result)=>{
      this.currentPackage = result.PackageName,
      this.currentActivity = result.ActivityName
    })


  },
  methods: {

  },
  beforeUnmount() {

  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>