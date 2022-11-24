from PySide6 import QtWidgets
import sys
from lib.exec import execute
import re
from PySide6.QtCore import (QCoreApplication, QDate, QDateTime, QLocale,
    QMetaObject, QObject, QPoint, QRect,
    QSize, QTime, QUrl, Qt)
from PySide6.QtGui import (QBrush, QColor, QConicalGradient, QCursor,
    QFont, QFontDatabase, QGradient, QIcon,
    QImage, QKeySequence, QLinearGradient, QPainter,
    QPalette, QPixmap, QRadialGradient, QTransform)
from PySide6.QtWidgets import (QAbstractButton, QApplication, QDialog, QDialogButtonBox,
                               QHeaderView, QSizePolicy, QTableView, QWidget, QListWidget, QListWidgetItem, QMessageBox,
                               QFormLayout, QToolButton, QMenu, QTableWidgetItem, QTableWidget, QFrame, QLabel)
from PySide6.QtCore import (QCoreApplication, QDate, QDateTime, QLocale,
    QMetaObject, QObject, QPoint, QRect,
    QSize, QTime, QUrl, Qt)
from PySide6.QtGui import (QBrush, QColor, QConicalGradient, QCursor,
    QFont, QFontDatabase, QGradient, QIcon,
    QImage, QKeySequence, QLinearGradient, QPainter,
    QPalette, QPixmap, QRadialGradient, QTransform)
from PySide6.QtWidgets import *


class AndroidDevice():
    def __init__(self):
        print(f"AndroidDevice init")
        self.devices_list = []
        self.current_device = {}
        self.get_all_device()

    def get_all_device(self)->None:
        self.devices_list = []
        ret = execute("adb devices -l")
        devices_info = [i for i in ret.split("\n") if 'model' in i]
        for info in devices_info:
            m = re.search(r'(.*?)device.*?model:([0-9A-Za-z]*)', info)
            serial = str(m.group(1).strip())  # 机器序列号
            model = str(m.group(2).strip())  # 机器型号
            self.devices_list.append({"serial": serial, "model": model})

    def get_all_attributes(self):
        device_prop = execute(f"adb -s {self.current_device['serial']} shell getprop")

        attitude = {"手机品牌":"ro.product.brand","设备名称":"ro.product.device","设备内部代号":"ro.product.model"}
        res = attitude.copy()
        for line in device_prop.splitlines():
            value = line.split(":")[-1].strip(" []")
            for k, v in attitude.items():
                if v in line:
                   res[k] = {"value":value,"source":v}
        self.current_device_attitude = res

class PageBase(QtWidgets.QMainWindow, AndroidDevice):
    def __init__(self):
        print("PageBase init")
        QtWidgets.QMainWindow.__init__(self)
        AndroidDevice.__init__(self)
        self.setWindowTitle("android工具箱 by saucerman")

    def create_button(self, QAbstractButton:QAbstractButton,object_name:str,qrect:QRect,text:str)->QPushButton:
        button = QPushButton(QAbstractButton)
        button.setObjectName(object_name)
        button.setGeometry(qrect)
        button.setText(text)
        return button

    # 所有继承类要实现下面2个接口
    def setupUi(self, MainWindow):
        pass
    def handel_event(self):
        pass

# 初始选择设备框框
class InitPage(PageBase):
    def __init__(self):
        super().__init__()


    def setupUi(self, MainWindow):
        MainWindow.resize(742, 542)
        # 在MainWindow创建butten和list todo: 这个列表页面太丑，待优化
        self.deviceListWidget = QListWidget(MainWindow)
        for device in self.devices_list:
            self.deviceListWidget.addItem(f"serial:{device['serial']}\tmodel:{device['model']}")
        self.deviceListWidget.setObjectName(u"deviceListWidget")
        self.deviceListWidget.setGeometry(QRect(210, 60, 321, 241))

        self.refreshButton = self.create_button(MainWindow, "refreshButton", QRect(220, 360, 111, 41), "刷新")
        self.enterButton = self.create_button(MainWindow, "enterButton", QRect(370, 360, 111, 41), "确定")

        self.handel_event()

    def handel_event(self):
        self.enterButton.clicked.connect(self.enter_click)
        self.refreshButton.clicked.connect(self.refresh_click)

    def refresh_click(self):
        self.get_all_device()
        print(f"refresh_click  devices_list:{self.devices_list}")
        self.deviceListWidget.clear()
        if not self.devices_list:
            return
        for device in self.devices_list:
            self.deviceListWidget.addItem(f"serial：{device['serial']}\tmodel:{device['model']}")

    def enter_click(self):
        # 默认情况下会返回第一个或者选中的那个，如果没有则说明没插入手机
        if not self.deviceListWidget.currentItem():
            QMessageBox.warning(self, '警告', '当前未插入手机')
        else:
            device_info = self.deviceListWidget.currentItem().text()
            self.device_serial = device_info.split("\t")[0].split(":")[1].strip()
            self.device_model = device_info.split("\t")[-1].split(":")[1].strip()
            self.current_device = {"serial":self.device_serial,"model":self.device_model}
            print(f"self.current_device:{self.current_device}")
            self.close()
            self.window = QtWidgets.QMainWindow()
            devicePage = DevicePage(self.current_device)
            devicePage.setupUi(self.window)
            self.window.show()

class DevicePage(PageBase):
    def __init__(self, current_device):
        super().__init__()
        print("DevicePage init")
        self.current_device = current_device
        self.get_all_attributes()

    def setupUi(self, MainWindow):
        if not MainWindow.objectName():
            MainWindow.setObjectName(u"MainWindow")
        MainWindow.resize(791, 512)

        self.centralwidget = QWidget(MainWindow)
        self.centralwidget.setObjectName(u"centralwidget")
        self.frame = QFrame(self.centralwidget)
        self.frame.setObjectName(u"frame")
        self.frame.setGeometry(QRect(100, 0, 701, 511))
        self.frame.setFrameShape(QFrame.StyledPanel)
        self.frame.setFrameShadow(QFrame.Raised)
        self.label = QLabel(self.frame)
        self.label.setObjectName(u"label")
        self.label.setText("设备信息")
        self.label.setGeometry(QRect(230, 20, 61, 31))

        self.tableWidget = QTableWidget(self.frame)
        self.tableWidget.setColumnCount(2)
        self.tableWidget.setColumnWidth(0,150)
        self.tableWidget.setColumnWidth(1,200)

        self.tableWidget.setHorizontalHeaderItem(0, QTableWidgetItem("值"))
        self.tableWidget.setHorizontalHeaderItem(1, QTableWidgetItem("来源"))
        self.tableWidget.setRowCount(len(self.current_device_attitude))
        i = 0
        for k, v in self.current_device_attitude.items():
            self.tableWidget.setVerticalHeaderItem(i, QTableWidgetItem(k))
            self.tableWidget.setItem(i, 0, QTableWidgetItem(v["value"]))
            self.tableWidget.setItem(i, 1, QTableWidgetItem(v["source"]))
            i += 1

        self.tableWidget.setObjectName(u"tableWidget")
        self.tableWidget.setGeometry(QRect(30, 60, 461, 261))
        self.tableWidget.horizontalHeader().setCascadingSectionResizes(False)
        self.homeButton = QPushButton(self.centralwidget)
        self.homeButton.setObjectName(u"homeButton")
        self.homeButton.setEnabled(True)
        self.homeButton.setText("主页")
        self.homeButton.setGeometry(QRect(0, 0, 100, 50))
        self.homeButton.setMinimumSize(QSize(100, 50))
        self.softButton = QPushButton(self.centralwidget)
        self.softButton.setObjectName(u"softButton")
        self.softButton.setText("软件信息")
        self.softButton.setGeometry(QRect(0, 50, 100, 50))
        self.softButton.setMinimumSize(QSize(100, 50))
        self.fileButton = QPushButton(self.centralwidget)
        self.fileButton.setObjectName(u"fileButton")
        self.fileButton.setText("文件管理")
        self.fileButton.setGeometry(QRect(0, 100, 100, 50))
        self.fileButton.setMinimumSize(QSize(100, 50))
        MainWindow.setCentralWidget(self.centralwidget)
        QMetaObject.connectSlotsByName(MainWindow)

class AndroidToolboxGui(PageBase):
    def __init__(self):
        print("AndroidToolboxGui init")
        self.initPage = InitPage()
    def run(self):
        self.initPage.show()


import sys

if __name__ == "__main__":
    app = QtWidgets.QApplication()
    app.setStyle('Fusion')
    initPage = InitPage()
    initPage.setupUi(initPage)
    initPage.show()
    sys.exit(app.exec_())
