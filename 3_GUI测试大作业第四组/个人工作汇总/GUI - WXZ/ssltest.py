import unittest
import re
import subprocess


class SSLTest(unittest.TestCase):

    def test_ssl(self):
        domain = 'zhiyuan.sjtu.edu.cn'
        cmd = 'curl -lvs https://{}/'.format(domain)
        sslinfo = subprocess.getstatusoutput(cmd)[1]

        print('domain:', domain)
        m = re.search(
            '(\*  subject.*?\*  SSL certificate.*?)\n', sslinfo, re.S)
        res = m.group(1)
        print(res)

        self.assertIn("SSL certificate verify ok", res)


if __name__ == '__main__':
    unittest.main()
