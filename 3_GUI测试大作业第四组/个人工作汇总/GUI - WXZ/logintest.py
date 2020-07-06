import unittest
from selenium import webdriver
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.support.ui import WebDriverWait


class LoginTest(unittest.TestCase):
    def setUp(self):
        self.driver = webdriver.Chrome()

    def test_login(self):
        driver = self.driver
        driver.maximize_window()
        driver.get('https://zhiyuan.sjtu.edu.cn/')
        self.assertIn('致远学院', driver.title)
        self.assertGreaterEqual(len(driver.get_cookies()), 2)
        logo = driver.find_element_by_xpath('//a[@href="/admin"]')
        WebDriverWait(driver, 10).until(
            lambda driver: logo.is_displayed())
        ActionChains(driver).click(logo).perform()
        self.assertRegex(r'^(上海交通大学统一身份认证|SJTU Single Sign On)$', driver.title)


def tearDown(self):
    self.driver.close()


if __name__ == '__main__':
    unittest.main()
