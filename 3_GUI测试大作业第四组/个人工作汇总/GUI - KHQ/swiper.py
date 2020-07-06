from selenium import webdriver
from time import sleep

def test_swiper():
    '''
        检测主页中的轮播图功能
    '''
    print('\n通过find_element_by_xpath方法进行测试...')
    try:
        driver = webdriver.Chrome('./chromedriver.exe')
        driver.get('https://zhiyuan.sjtu.edu.cn/html/zhiyuan/index.php')

        pictures = set()
        for i in range(30):
            pictures.add(driver.find_element_by_class_name('swiper-slide-active').find_element_by_tag_name('img').get_attribute('src'))
            print('等待轮播图播放...(%ds)' % (30-i), end='\r')
            sleep(1)
        assert pictures.__len__() == 4
        paths = ['https://zhiyuan.sjtu.edu.cn/images/articleImages/original/php0bBTJb1589185690.png',
                 'https://zhiyuan.sjtu.edu.cn/images/articleImages/original/phpYq3AR11589186116.png',
                 'https://zhiyuan.sjtu.edu.cn/images/articleImages/original/phpasUtIl1585035323.jpg',
                 'https://zhiyuan.sjtu.edu.cn/images/articleImages/original/phpo1kIof1585035323.jpg']
        for path in paths:
            assert path in pictures
        print('轮播图正常自动播放完毕')
        
        for i in range(1, 5):
            driver.find_element_by_xpath('/html/body/div[1]/div/div[2]/span[%d]' % i).click()
            assert 'swiper-pagination-bullet-active' in driver.find_element_by_xpath('/html/body/div[1]/div/div[2]/span[%d]' % i).get_attribute('class')
        print('测试点击控制轮播图完毕')

    finally:
        sleep(3)
        driver.quit()


if __name__ == '__main__':
    test_swiper()
    print('\n所有测试全部通过！')
