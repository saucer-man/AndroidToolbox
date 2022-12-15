<template>
  <el-form
    ref="ruleFormRef"
    :model="adbForm"
    :rules="rules"
    label-width="120px"
    status-icon
  >
    <el-form-item label="ADB命令" prop="command" size="large">
      <el-input v-model="adbForm.command" />
    </el-form-item>

    <el-form-item label="输出" prop="output">
      <el-input v-model="adbForm.output" :rows="17" type="textarea" />
    </el-form-item>

    <div style="text-align: center">
      <el-button type="primary" @click="submitForm(ruleFormRef)"
        >提交</el-button
      >
      <el-button @click="resetForm(ruleFormRef)">重置</el-button>
    </div>
  </el-form>
</template>

<script lang="ts" setup>
import { reactive, ref, watch } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { ListPath, Excute } from "../../wailsjs/go/app/App.js";
import { ElMessage, ElMessageBox } from "element-plus";
const props = defineProps({
  selectdevice: {
    type: String,
    required: true,
  },
});
watch(
  () => props.selectdevice,
  (newProps) => {
    console.log("通过prop获取到的selectdevice更新成了:", props.selectdevice);
    adbForm.command = `adb -s ${props.selectdevice} shell ls /storage/emulated/0`;
  }
);
const ruleFormRef = ref<FormInstance>();
const adbForm = reactive({
  command: `adb -s ${props.selectdevice} shell ls /storage/emulated/0`,
  output: "",
});
const validateCommand = (rule: any, value: any, callback: any) => {
  if (value === "") {
    callback(new Error("请输入命令"));
  } else {
    if (!adbForm.command.startsWith("adb")) {
      callback(new Error("命令必须以adb开头"));
    }
    callback();
  }
};
const rules = reactive<FormRules>({
  name: [{ validator: validateCommand, trigger: "blur" }],
});

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid) => {
    if (valid) {
      Excute(adbForm.command.split(" ")).then((result) => {
        console.log("handleCommandResult返回值:", result);
        if (result.ExitCode == 0) {
          // ElMessage.success("命令执行成功")
          adbForm.output = result.Stdout;
        } else {
          // ElMessage.success("命令执行出错: " + result.Stderr)
          adbForm.output = result.Stderr;
        }
      });
    } else {
      console.log("error submit!");
      return false;
    }
  });
};

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields();
};
</script>
<style type="text/css">
.el-input {
  width: 600px;
}

.el-form-item {
  text-align: center;
}
</style>
