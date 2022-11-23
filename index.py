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
                               QFormLayout, QToolButton, QMenu)
from PySide6.QtCore import (QCoreApplication, QDate, QDateTime, QLocale,
    QMetaObject, QObject, QPoint, QRect,
    QSize, QTime, QUrl, Qt)
from PySide6.QtGui import (QBrush, QColor, QConicalGradient, QCursor,
    QFont, QFontDatabase, QGradient, QIcon,
    QImage, QKeySequence, QLinearGradient, QPainter,
    QPalette, QPixmap, QRadialGradient, QTransform)
from PySide6.QtWidgets import (QApplication, QListView, QMainWindow, QMenuBar,
    QPushButton, QSizePolicy, QStatusBar, QWidget)


class AndroidDevice():
    def __init__(self):
        print(f"AndroidDevice init")
        self.devices_list = []

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

class PageBase(QtWidgets.QMainWindow, AndroidDevice):
    def __init__(self):
        print("AndroidToolboxGui init")
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
        self.setupUi(self)
        self.handel_event()

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
            print(self.deviceListWidget.currentItem().text()) # serial：a464ccdf0410	model:M2004J19C

            self.close()
            devicePage = DevicePage(self.deviceListWidget.currentItem().text())
            devicePage.show()

class DevicePage(PageBase):
    def __init__(self, device_info: str):
        super().__init__()
        print("DevicePage init")
        self.device_serial = device_info.split("\t")[0].split(":")[1].strip()
        self.device_model = device_info.split("\t")[-1].split(":")[1].strip()
        print(f"得到的参数为serial:{self.device_serial}\tmodel:{self.device_model}")

        self.setupUi(self)
        self.handel_event()

    def setupUi(self, MainWindow):
        if not MainWindow.objectName():
            MainWindow.setObjectName(u"MainWindow")
        MainWindow.resize(800, 600)
        self.centralwidget = QWidget(MainWindow)
        self.centralwidget.setObjectName(u"centralwidget")
        self.listView = QListView(self.centralwidget)
        self.listView.setObjectName(u"listView")
        self.listView.setGeometry(QRect(0, -20, 81, 571))
        self.formLayoutWidget = QWidget(self.centralwidget)
        self.formLayoutWidget.setObjectName(u"formLayoutWidget")
        self.formLayoutWidget.setGeometry(QRect(79, -11, 721, 561))
        self.formLayout = QFormLayout(self.formLayoutWidget)
        self.formLayout.setObjectName(u"formLayout")
        self.formLayout.setContentsMargins(0, 0, 0, 0)
        self.toolButton = QToolButton(self.formLayoutWidget)
        self.toolButton.setObjectName(u"toolButton")

        self.formLayout.setWidget(0, QFormLayout.LabelRole, self.toolButton)

        MainWindow.setCentralWidget(self.centralwidget)
        self.menubar = QMenuBar(MainWindow)
        self.menubar.setObjectName(u"menubar")
        self.menubar.setGeometry(QRect(0, 0, 800, 24))
        self.menu = QMenu(self.menubar)
        self.menu.setObjectName(u"menu")
        MainWindow.setMenuBar(self.menubar)
        self.statusbar = QStatusBar(MainWindow)
        self.statusbar.setObjectName(u"statusbar")
        MainWindow.setStatusBar(self.statusbar)

        self.menubar.addAction(self.menu.menuAction())

        self.retranslateUi(MainWindow)

        QMetaObject.connectSlotsByName(MainWindow)
    # setupUi

    def retranslateUi(self, MainWindow):
        MainWindow.setWindowTitle(QCoreApplication.translate("MainWindow", u"MainWindow", None))
        self.toolButton.setText(QCoreApplication.translate("MainWindow", u"\u4f60\u597d", None))
        self.menu.setTitle(QCoreApplication.translate("MainWindow", u"\u4f60\u597d", None))
    # retranslateUi

    def handel_event(self):
        pass



class AndroidToolboxGui(PageBase):
    def __init__(self):
        print("AndroidToolboxGui init")
        self.initPage = InitPage()
    def run(self):
        self.initPage.show()


import sys

if __name__ == "__main__":
    app = QtWidgets.QApplication()
    w = InitPage()
    w.show()

    sys.exit(app.exec_())
