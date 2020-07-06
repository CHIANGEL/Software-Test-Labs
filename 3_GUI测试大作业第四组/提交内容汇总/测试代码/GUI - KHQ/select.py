from selenium import webdriver
from time import sleep

def test_select():
    '''
        检测课程信息页面下的下拉框和筛选功能
    '''
    print('\n通过find_element_by_xpath方法进行测试...')
    try:
        driver = webdriver.Chrome('./chromedriver.exe')
        driver.get('https://zhiyuan.sjtu.edu.cn/html/zhiyuan/course_list.php?bg=jxfw')

        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[1]/div/div/input').click()
        assert [ele.text for ele in driver.find_elements_by_xpath('/html/body/div[2]/div/form/div/div[1]/div/div[1]/div/dl/dd')] == [
            '学年', '2019-2020学年', '2018-2019学年', '2017-2018学年', '2016-2017学年', '2015-2016学年', '2014-2015学年', '2013-2014学年', '2012-2013学年', '2011-2012学年', '2010-2011学年', '2009-2010学年', '2008-2009学年']
        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[1]/div/dl/dd[3]').click()  # 2018-2019学年
        assert driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[1]/div/dl/dd[3]').get_attribute('class') == 'layui-this'
        print('选择2018-2019学年')


        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[2]/div/div/input').click()
        assert [ele.text for ele in driver.find_elements_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[2]/div/dl/dd')] == ['学期', '秋季', '春季', '夏季']
        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[2]/div/dl/dd[2]').click()  # 秋季
        assert driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[2]/div/dl/dd[2]').get_attribute('class') == 'layui-this'
        print('选择秋季')

        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[3]/div/div/input').click()
        assert [ele.text for ele in driver.find_elements_by_xpath('/html/body/div[2]/div/form/div/div[1]/div/div[3]/div/dl/dd')] == [
            '专业', '数理科学', '数学', '物理学', '生命科学', '计算机科学', '化学', '工科', '生物医学科学']
        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[3]/div/dl/dd[6]').click()  # 计算机科学
        assert driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[3]/div/dl/dd[6]').get_attribute('class') == 'layui-this'
        print('选择计算机科学')

        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[4]/div/div/input').click()
        assert [ele.text for ele in driver.find_elements_by_xpath('/html/body/div[2]/div/form/div/div[1]/div/div[4]/div/dl/dd')] == [
            '年级', '2019级', '2018级', '2017级', '2016级', '2015级', '2014级', '2013级', '2012级', '2011级', '2010级', '2009级', '2008级']
        driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[4]/div/dl/dd[4]').click()  # 2017级
        assert driver.find_element_by_xpath(
            '/html/body/div[2]/div/form/div/div[1]/div/div[4]/div/dl/dd[4]').get_attribute('class') == 'layui-this'
        print('选择2017级')

        driver.find_element_by_xpath('/html/body/div[2]/div/form/div/div[2]/button').click() # 筛选
        print('筛选')
        sleep(3) # 等待结果返回

        assert driver.find_elements_by_class_name('course-item').__len__() == 9
        for ele in driver.find_elements_by_class_name('course-item'):
            assert '2018 秋季 计算机科学 2017' in ele.text
        print('筛选结果正确')

    finally:
        sleep(3)
        driver.quit()


if __name__ == '__main__':
    test_select()
    print('\n所有测试全部通过！')
