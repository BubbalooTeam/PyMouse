#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from ... import UtilsDecorator
from speedtest import Speedtest

class NetworkUtils:
    @staticmethod
    @UtilsDecorator().aiowrap()
    def speedtest_performer():
        run = Speedtest()
        bs = run.get_best_server()
        dl = round(run.download() / 1024 / 1024, 2)
        ul = round(run.upload() / 1024 / 1024, 2)
        # //
        run.results.share()
        result = run.results.dict()
        name = result.get("server", {}).get("name", "No Infos")
        host = bs.get("sponsor", "No Infos")
        ping = bs.get("latency", "No Infos")
        isp = result.get("client", {}).get("isp", "No Infos")
        country = result.get("server", {}).get("country", "No Infos")
        cc = result.get("server", {}).get("cc", "No Infos")
        path = (result.get("share", NetworkEvents.FAIL_IN_EXTRACT_PHOTOS))
        return dl, ul, name, host, ping, isp, country, cc, path
    
class NetworkEvents:
    RUNNING_SPEEDTEST = "https://telegra.ph/file/4beef2f3139e3b160615f.jpg"
    FAIL_IN_EXTRACT_PHOTOS = "https://telegra.ph/file/da2852745e2cc6160e40e.jpg"