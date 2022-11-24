import os
import subprocess


def execute(command):
    """
    同步执行，确保会得到返回值
    """
    try:
        proc = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        outs, errs = proc.communicate(timeout=1)
    except subprocess.TimeoutExpired:
        proc.kill()
        outs, errs = proc.communicate()
    if errs:
        print(f"执行{command}出错了:{errs}")
        return ""
    return outs.decode()


if __name__=="__main__":
    device_prop = execute(f"adb shell getprop")

    attitude = {"手机品牌":"ro.product.brand","设备名称":"ro.product.device","设备内部代号":"ro.product.model"}
    res = attitude.copy()
    for line in device_prop.splitlines():
        value = line.split(":")[-1].strip(" []")
        for k, v in attitude.items():
            if v in line:
               res[k] = {"value":value,"source":v}
    print(res)

