## LJH 部分测试简介

## Tools
pytest 作为测试框架，负责断言
selenium 配合 driver 进行相应的操作

## NOTE
1. selenium一共有八种定位元素的基本方法，我尽量进行多种运用，对不同的按钮、输入框用不同的定位方法进行测试。
2. 每次有页面跳转或者导航行为后，手动sleep三秒时间，防止网络延迟问题带来的执行流bug

## 测试内容
- 按钮测试：采用三种定位方法对三类按钮进行测试
  1. test_button_by_xpath()
  2. test_button_by_class_name()
  3. test_button_by_link_text()
- 输入框测试：全站一共有两种输入框，一个是在全站搜索页面的输入框，一个是在主页点击搜索图标后弹出的输入框，也采用三种定位方法进行测试
  1. test_inputBox_by_id()
  2. test_inputBox_by_tag_name()
  3. test_inputBox_by_CSS_Selectors()