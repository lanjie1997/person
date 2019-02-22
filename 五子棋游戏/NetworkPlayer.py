# -*- coding:utf-8 -*-
__author__ = 'lanjie'
__Date__ = '2018/8/30 9:55'

from PyQt5.QtWidgets import QPushButton,QLabel

from Base import BasePlayer,TDPushButton


class NetworkPlayer(BasePlayer):
    def __init__(self,parent=None):
        super().__init__(parent)
        self.setup_ui()

    def setup_ui(self):
        super().setup_ui()
        self.state_label = QLabel("游戏状态",self)
        self.state_text = QLabel("等待连接",self)
        self.state_label.move(630,200)
        self.state_text.move(680,204)

        self.cuicu_btn = TDPushButton("source/催促按钮_normal.png","source/催促按钮_hover.png","source/催促按钮_press.png",self)
        self.cuicu_btn.show()
        self.cuicu_btn.move(640,450)





