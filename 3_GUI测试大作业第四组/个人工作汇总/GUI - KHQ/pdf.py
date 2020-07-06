from selenium import webdriver
import time
from selenium.webdriver.common.keys import Keys
import win32gui, win32con
import re

def test_pdf():
    '''
        检测教学服务-培养方案中pdf功能
    '''
    print('\n通过find_element_by_xpath方法进行测试...')
    try:
        driver = webdriver.Chrome('./chromedriver.exe')
        driver.get('https://zhiyuan.sjtu.edu.cn/html/zhiyuan/service_list.php?bg=jxfw')

        time.sleep(1)
        driver.find_element_by_xpath('/html/body/div[2]/div/form/div/div[1]/div/div/input').click()
        driver.find_element_by_xpath('/html/body/div[2]/div/form/div/div[1]/div/dl/dd[6]').click()
        driver.find_element_by_xpath('/html/body/div[2]/div/form/div/div[2]/div/div').click()
        driver.find_element_by_xpath('/html/body/div[2]/div/form/div/div[2]/div/dl/dd[5]').click()
        driver.find_element_by_xpath('/html/body/div[2]/div/form/div/div[3]/button').click()
        print('查找pdf')

        time.sleep(3)
        iframe = driver.find_element_by_xpath('/html/body/div[3]/div/div/iframe')
        print(iframe.get_attribute('src'))
        assert re.match('https://zhiyuan.sjtu.edu.cn/DevelopmentPlan/.*.pdf', iframe.get_attribute('src'))

        iframe.click()
        iframe.send_keys(Keys.TAB)
        iframe.send_keys(Keys.TAB)
        time.sleep(1)
        iframe.send_keys(Keys.ENTER)
        iframe.send_keys(Keys.ENTER)
        iframe.send_keys(Keys.ENTER)
        iframe.send_keys(Keys.ENTER)
        print('测试旋转完成')

        time.sleep(1)
        iframe.send_keys(Keys.TAB)
        iframe.send_keys(Keys.ENTER)
        time.sleep(3)
        dialog = win32gui.FindWindow('#32770', u'另存为')
        button = win32gui.FindWindowEx(dialog, 0, 'Button', None)  # 确定按钮Button
        time.sleep(3)
        assert 0 == win32gui.SendMessage(dialog, win32con.WM_COMMAND, 1, button)  # 按button
        print('测试下载完成')

        driver.switch_to.frame(iframe)
        assert driver.find_element_by_xpath('/html/body/embed').get_attribute('type') == 'application/pdf'

    finally:
        time.sleep(3)
        driver.quit()


if __name__ == '__main__':
    test_pdf()
    print('\n所有测试全部通过！')
