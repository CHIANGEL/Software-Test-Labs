from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver import ActionChains
from time import sleep

def test_button_by_xpath():
    '''
     通过find_element_by_xpath的方法，检测【主页】的三个“查看更多”按钮
    '''
    print('\n通过find_element_by_xpath方法进行测试...')
    print('\n==>检测第一个按钮元素：主页->最新动态->查看更多')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn')
    sleep(3)
    block_index_news_div = driver.find_element_by_xpath('/html/body/div[2]')
    print(' - class name of div: {}'.format(block_index_news_div.get_attribute('class')))
    assert block_index_news_div.get_attribute('class') == 'block index-news'
    button = driver.find_element_by_xpath('/html/body/div[2]/div/div/a')
    print(' - class name of button: {}'.format(button.get_attribute('class')))
    print(' - href of button: {}'.format(button.get_attribute('href')))
    assert button.get_attribute('class') == 'more-btn btn-yellow'
    assert button.get_attribute('textContent') == '查看更多'
    button.click()
    sleep(3)
    print(' - page title after click: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '最新动态 - 上海交通大学致远学院'
    driver.quit()

    print('\n==>检测第二个按钮元素：主页->重要活动->查看更多')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn')
    sleep(3)
    block_div = driver.find_element_by_xpath('/html/body/div[3]')
    print(' - class name of div: {}'.format(block_div.get_attribute('class')))
    assert block_div.get_attribute('class') == 'block grey'
    button = driver.find_element_by_xpath('/html/body/div[3]/div/div/a')
    print(' - class name of button: {}'.format(button.get_attribute('class')))
    print(' - href of button: {}'.format(button.get_attribute('href')))
    assert button.get_attribute('class') == 'more-btn btn-yellow'
    assert button.get_attribute('textContent') == '查看更多'
    button.click()
    sleep(3)
    print(' - page title after click: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '重要活动 - 上海交通大学致远学院'
    driver.quit()

    print('\n==>检测第三个按钮元素：主页->重要通知->查看更多')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn')
    sleep(3)
    block_div = driver.find_element_by_xpath('/html/body/div[4]')
    print(' - class name of div: {}'.format(block_div.get_attribute('class')))
    assert block_div.get_attribute('class') == 'block'
    button = driver.find_element_by_xpath('/html/body/div[4]/div/div/a')
    print(' - class name of button: {}'.format(button.get_attribute('class')))
    print(' - href of button: {}'.format(button.get_attribute('href')))
    assert button.get_attribute('class') == 'more-btn btn-yellow'
    assert button.get_attribute('textContent') == '查看更多'
    button.click()
    sleep(3)
    print(' - page title after click: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '重要通知 - 上海交通大学致远学院'

    print('\n==>find_element_by_xpath方法检测结束')
    driver.quit()

def test_button_by_class_name():
    '''
     通过find_element_by_class_name方法测试【主页->最新动态】的活动按钮与页码跳转按钮
    '''
    def test_slides(driver):
        slide_divs = driver.find_elements_by_class_name('slide')
        silde_num = len(slide_divs)
        print(' - Number of the slides: {}'.format(silde_num))
        for idx in range(silde_num):
            original_title = (driver.find_elements_by_class_name('title'))[idx].get_attribute('textContent').strip()
            (driver.find_elements_by_class_name('slide'))[idx].click()
            sleep(3)
            current_title = (driver.find_element_by_class_name('news-title')).text.strip()
            print(' - title {} before click: {}'.format(idx + 1, original_title))
            print(' - title {} after  click: {}'.format(idx + 1, current_title))
            assert original_title == current_title
            driver.back()
            sleep(3)
    
    print('\n通过find_element_by_class_name方法进行测试...')
    print('\n==>检测活动按钮：主页->最新动态->活动按钮')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn/html/zhiyuan/news_list.php')
    sleep(3)
    print(' - current page title: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '最新动态 - 上海交通大学致远学院'
    test_slides(driver)
    
    print('\n==>检测页码跳转按钮：主页->最新动态->页码跳转按钮')
    target_page = 5
    ((driver.find_element_by_class_name('pages')).find_elements_by_tag_name('a'))[target_page - 1].click()
    sleep(3)
    current_page = int(((driver.find_element_by_class_name('layui-laypage-curr')).find_elements_by_tag_name('em'))[-1].get_attribute('textContent'))
    print(' - target  page: {}'.format(target_page))
    print(' - current page: {}'.format(current_page))
    assert target_page == current_page

    print('\n==>重新检测页码跳转后的活动按钮：主页->最新动态->页码跳转按钮->活动按钮')
    print(' - current page title: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '最新动态 - 上海交通大学致远学院'
    test_slides(driver)

    print('\n==>find_element_by_class_name方法检测结束')
    driver.quit()

def test_button_by_link_text():
    '''
     通过find_element_by_link_text方法测试【主页->导航栏】的导航按钮
    '''
    print('\n通过find_element_by_link_text方法进行测试...')
    print('\n==>检测主页的各类导航按钮')
    driver = webdriver.Chrome()
    driver.get('https://zhiyuan.sjtu.edu.cn/')
    actionChains = ActionChains(driver)
    sleep(3)
    print(' - current page title: {}'.format(driver.find_element_by_xpath('/html/head/title').get_attribute('textContent')))
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '首页 - 上海交通大学致远学院 - sh'

    student_service = driver.find_element_by_link_text('学生工作')
    actionChains.move_to_element(student_service).perform()
    sleep(3)
    (driver.find_element_by_link_text('学生服务')).click()
    sleep(3)
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '重要通知 - 上海交通大学致远学院'
    print(' - current pagesnavigatio: {}'.format('主页->学生工作->学生服务->重要通知'))
    
    link_text = '奖助学金'
    print(' - target navigation: {}'.format(link_text))
    (driver.find_element_by_link_text(link_text)).click()
    sleep(3)
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '重要通知 - 上海交通大学致远学院'
    assert (driver.find_element_by_link_text(link_text)).get_attribute('class') == 'active'
    
    link_text = '学生事务'
    print(' - target navigation: {}'.format(link_text))
    (driver.find_element_by_link_text(link_text)).click()
    sleep(3)
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '重要通知 - 上海交通大学致远学院'
    assert (driver.find_element_by_link_text(link_text)).get_attribute('class') == 'active'
    
    link_text = '团学通知'
    print(' - target navigation: {}'.format(link_text))
    (driver.find_element_by_link_text(link_text)).click()
    sleep(3)
    assert driver.find_element_by_xpath('/html/head/title').get_attribute('textContent') == '重要通知 - 上海交通大学致远学院'
    assert (driver.find_element_by_link_text(link_text)).get_attribute('class') == 'active'

    print('\n==>find_element_by_link_text方法检测结束')
    driver.quit()

if __name__ == '__main__':
    test_button_by_xpath()
    test_button_by_class_name()
    test_button_by_link_text()
    print('\n所有测试全部通过！')