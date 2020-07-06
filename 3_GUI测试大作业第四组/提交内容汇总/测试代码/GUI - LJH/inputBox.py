from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from time import sleep

def test_inputBox_by_id():
    '''
     通过find_element_by_id的方法，检测【主页->整站搜索】的输入框
    '''
    print('\n==>通过find_element_by_id方法进行测试...')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn/html/zhiyuan/search.php')
    sleep(3)
    print(' - current page title: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '整站搜索 - 上海交通大学致远学院'
    
    try:
        laypage = driver.find_element_by_id('layui-laypage-1')
        print(' - ERROR! Should not have lay pages here!')
        assert 0
    except:
        print(' - currently no lay pages here')

    search_text = '致远沙龙'
    inputBox = driver.find_element_by_id('keyword')
    inputBox.send_keys(search_text)
    inputBox.send_keys(Keys.ENTER)
    sleep(3)

    try:
        laypage = driver.find_element_by_id('layui-laypage-1')
        print(' - find lay pages after searching')
    except:
        print(' - ERROR! Should have lay pages after searching!')
    
    driver.quit()

def test_inputBox_by_tag_name():
    '''
     通过find_element_by_tag_name的方法，检测【主页->整站搜索】的输入框
    '''
    print('\n==>通过find_element_by_tag_name方法进行测试...')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn/html/zhiyuan/search.php')
    sleep(3)
    print(' - current page title: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '整站搜索 - 上海交通大学致远学院'
    
    try:
        laypage = driver.find_element_by_id('layui-laypage-1')
        print(' - ERROR! Should not have lay pages here!')
        assert 0
    except:
        print(' - currently no lay pages here')

    search_text = '致远沙龙'
    inputBox = driver.find_element_by_tag_name('input')
    inputBox.send_keys(search_text)
    inputBox.send_keys(Keys.ENTER)
    sleep(3)

    try:
        laypage = driver.find_element_by_id('layui-laypage-1')
        print(' - find lay pages after searching')
    except:
        print(' - ERROR! Should have lay pages after searching!')
    
    driver.quit()

def test_inputBox_by_CSS_Selectors():
    '''
     通过find_element_by_CSS_Selectors的方法，检测【主页->整站搜索】的输入框
    '''
    print('\n==>通过find_element_by_CSS_Selectors方法进行测试...')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn/')
    sleep(3)
    print(' - current page title: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '首页 - 上海交通大学致远学院 - sh'
    
    try:
        laypage = driver.find_element_by_id('layui-laypage-1')
        print(' - ERROR! Should not have lay pages here!')
        assert 0
    except:
        print(' - currently no lay pages here')

    (driver.find_element_by_css_selector('.search-btn')).click()
    sleep(3)
    search_text = '奖学金'
    inputBox = driver.find_element_by_css_selector('#search-input')
    inputBox.send_keys(search_text)
    inputBox.send_keys(Keys.ENTER)
    sleep(3)

    try:
        laypage = driver.find_element_by_id('layui-laypage-1')
        print(' - find lay pages after searching')
    except:
        print(' - ERROR! Should have lay pages after searching!')
    
    driver.quit()

if __name__ == '__main__':
    test_inputBox_by_id()
    test_inputBox_by_tag_name()
    test_inputBox_by_CSS_Selectors()
    print('\n所有测试全部通过！')