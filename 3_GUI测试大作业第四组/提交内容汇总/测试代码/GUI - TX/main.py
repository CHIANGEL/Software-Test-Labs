from selenium import webdriver
import time
import requests as req
from PIL import Image
from io import BytesIO
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.chrome.options import Options
import os
from time import sleep


class TestImagePixs(object):

    def setup_method(self, method):
        self.driver = webdriver.Chrome(executable_path="./chromedriver.exe")

    def teardown_method(self, method):
        self.driver.quit()

    def loadImages(self, pics):
        pixs = []
        for item in pics:
            res = req.get(item)
            img = Image.open(BytesIO(res.content))
            pixs.append(img.size)
        return pixs

    def test_pixs(self):
        pics = []
        self.driver.get('https://zhiyuan.sjtu.edu.cn')
        pics.append(self.driver.find_element_by_xpath(
            '/html/body/div[1]/div/div[1]/div[3]/img').get_attribute("src"))
        pics.append(self.driver.find_element_by_xpath(
            '/html/body/div[1]/div/div[1]/div[4]/img').get_attribute("src"))
        pics.append(self.driver.find_element_by_xpath(
            '/html/body/div[1]/div/div[1]/div[5]/img').get_attribute("src"))
        pics.append(self.driver.find_element_by_xpath(
            '/html/body/div[1]/div/div[1]/div[6]/img').get_attribute("src"))
        pixs = self.loadImages(pics)
        for item in pixs:
            assert item[0] == 1920 and item[1] == 848

    def test_calendar(self):
        self.driver.get('https://zhiyuan.sjtu.edu.cn')
        self.driver.find_element_by_xpath(
            '//*[@id="layui-laydate1"]/div/div[1]/div/span[2]').click()  # click month selector
        self.driver.find_element_by_xpath(
            '//*[@id="layui-laydate1"]/div/div[2]/ul/li[5]').click()  # choose month
        self.driver.find_element_by_xpath(
            '//*[@id="layui-laydate1"]/div/div[2]/table/tbody/tr[3]/td[3]').click()  # click a invalid date
        assert self.driver.current_url == 'https://zhiyuan.sjtu.edu.cn/html/zhiyuan/'
        self.driver.find_element_by_xpath(
            '//*[@id="layui-laydate1"]/div/div[2]/table/tbody/tr[2]/td[4]').click()  # click a valid date
        assert self.driver.current_url == 'https://zhiyuan.sjtu.edu.cn/html/zhiyuan/events_list.php?date=2020-05-06'

    def test_download(self):
        os.remove("./{2018级}-{生命科学方向}.pdf")
        options = Options()
        options.add_experimental_option("prefs", {
            "download.default_directory": os.getcwd(),
            "download.prompt_for_download": False,
            "download.directory_upgrade": True,
            "safebrowsing.enabled": True,
            "plugins.always_open_pdf_externally": True
        })

        try:
            driver = webdriver.Chrome(
                executable_path="./chromedriver.exe", chrome_options=options)
            driver.maximize_window()
            driver.get('https://zhiyuan.sjtu.edu.cn')
            root_action = ActionChains(driver)
            element = driver.find_element_by_xpath(
                '//*[@id="nav-animate"]/li[6]/a')
            root_action.move_to_element(element).perform()
            redirect_action = ActionChains(driver)
            redirect_action.click(driver.find_element_by_xpath(
                '//*[@id="nav-animate"]/li[6]/ul/div/div/li[2]/a'))
            redirect_action.perform()
            driver.find_element_by_xpath(
                '/html/body/div[2]/div/form/div/div[1]/div/div/input').click()  # click major select
            driver.find_element_by_xpath(
                '/html/body/div[2]/div/form/div/div[1]/div/dl/dd[5]').click()  # select major
            driver.find_element_by_xpath(
                '/html/body/div[2]/div/form/div/div[2]/div/div').click()  # click grade select
            driver.find_element_by_xpath(
                '/html/body/div[2]/div/form/div/div[2]/div/dl/dd[4]').click()  # select grade
            driver.find_element_by_xpath(
                '/html/body/div[2]/div/form/div/div[3]/button').click()
            driver.implicitly_wait(10)
            frame = driver.find_element_by_xpath(
                '/html/body/div[3]/div/div/iframe')
            pdf_src = frame.get_attribute("src")
            assert pdf_src != ''
            res = req.get(pdf_src)
            assert len(res.content) != 0
            driver.switch_to.frame(frame)
            driver.find_element_by_xpath('//*[@id="main-content"]/a').click()
        finally:
            sleep(10) # wait for file download
            driver.quit()
        assert os.path.exists(os.getcwd() + "/{2018级}-{生命科学方向}.pdf")
