#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

# Some parts of this code are based on the PyKorone project
# (https://github.com/HitaloM/PyKorone/tree/main/src/korone/modules/gsm_arena), which is licensed under the 3-Clause BSD license.

from bs4 import BeautifulSoup
from dataclasses import dataclass
from httpx import HTTPError
from typing import Any

from pymouse.utils import http
from pymouse.utils.tools.gsm_arena.exceptions import (
    GSMarenaBadRequest,
    GSMarenaCategoryError,
    GSMarenaDeviceNotFound,
    GSMarenaManyRequests,
)

@dataclass(frozen=True, slots=True)
class GSMarenaDeviceSearchResult:
    id: str | None
    name: str | None
    image: str | None
    description: str | None

@dataclass(frozen=True, slots=True)
class GSMarenaSearchResults:
    results: list[GSMarenaDeviceSearchResult]

@dataclass(frozen=True, slots=True)
class GSMarenaDeviceBaseResult:
    name: str | None
    image: str | None
    phone_details: list[dict]

@dataclass(frozen=True, slots=True)
class GSMarenaBaseFormatResult:
    category: str
    attrs: list[str]

class GSMarena:
    def __init__(self):
        self.GSMarenaByPassURL: str = "https://cors-bypass.amano.workers.dev/https://www.gsmarena.com/{url}"
        self.Headers: dict = {
            "accept-language": "en-US,en;q=0.9",
            "cache-control": "max-age=0",
            "priority": "u=0, i",
            "user-agent": "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)",
            "referer": "https://www.gsmarena.com"
        }

    async def getDataFromURL(self, url: str) -> str:
        base_url = self.GSMarenaByPassURL.format(
            url=url
        )
        r = await http.get(
            url=base_url,
            headers=self.Headers,
        )
        if r.status_code != 200:
            if r.status_code == 400:
                msg = "An error occurred in the Client (BadRequest), please try again."
                raise GSMarenaBadRequest(msg)
            if r.status_code == 429:
                msg = "Please wait while we automatically unlock PyMouse access. This process may take minutes, hours, days, or even weeks."
                raise GSMarenaManyRequests(msg)
            else:
                msg = "[GSMarena]: An unknown error occurred while requesting on the GSMarena website."
                raise HTTPError(msg)
        else:
            return r.text

    async def search_device(self, query: str) -> GSMarenaSearchResults:
        html = await self.getDataFromURL(
            url="results.php3?sQuickSearch=yes&sName={query}".format(
                query=query,
            )
        )
        soup = BeautifulSoup(html, "html.parser")

        GSMarenaList = []

        devices = soup.select(".makers li")

        for device in devices:
            imageBlock = device.find("img")
            GSMarenaList.append(
                GSMarenaDeviceSearchResult(
                    id=device.find("a").get("href").replace(".php", ""),
                    name=device.find("span").text.replace("\n", "").replace("\r", "").replace("\t", ""),
                    image=imageBlock.get("src"),
                    description=imageBlock.get("title")
                ),
            )
        if GSMarenaList:
            return GSMarenaSearchResults(
                results=GSMarenaList,
            )
        else:
            msg = "The device was not found."
            raise GSMarenaDeviceNotFound(msg)

    async def fetch_device(self, device: str) -> GSMarenaDeviceBaseResult:
        html = await self.getDataFromURL(
            url="{device}.php".format(
                device=device,
            ),
        )
        soup = BeautifulSoup(html, 'html.parser')

        name_base = soup.find(class_='specs-phone-name-title')
        name = name_base.get_text() if name_base else "N/A"
        image = soup.find('div', class_='specs-photo-main').find('img').get('src')
        tableSpec = soup.find_all('table')
        phoneDetails = []

        for tS in tableSpec:
            phoneSpecifications = []
            category_base = tS.find('th')
            category = category_base.get_text() if category_base else "N/A"
            specification_n = tS.find_all('tr')

            for sN in specification_n:
                ttl = sN.find('td', class_='ttl')
                nfo = sN.find('td', class_='nfo')
                if ttl and nfo:
                    phoneSpecifications.append({
                        'name': ttl.get_text(),
                        'value': nfo.get_text(),
                    })
            if category != "N/A":
                phoneDetails.append({
                    "category": category,
                    "specifications": phoneSpecifications,
                })
            else:
                pass
        return GSMarenaDeviceBaseResult(
            name=name,
            image=image,
            phone_details=phoneDetails,
        )

    @staticmethod
    def get_specifications(specifications: list[dict[str, Any]], category: str, attrs: list[str]) -> str:
        categoryData = next((item for item in specifications if item["category"] == category), None)

        if not categoryData:
            msg = "The {category} category cannot be found, please refine your code.".format(
                category=category
            )
            raise GSMarenaCategoryError(msg)

        if attrs:
            return "\n".join(spec.get("value")
                     for spec in categoryData.get("specifications")
                     if spec.get("name") in attrs)
        else:
            return "\n".join(spec.get("value")
                             for spec in categoryData["specifications"])

    def parse_specifications(self, specifications: list[dict[str, Any]]) -> dict:
        specifications_map = {
            "status": GSMarenaBaseFormatResult(category="Launch", attrs=["Status"]),
            "network": GSMarenaBaseFormatResult(category="Network", attrs=["Technology"]),
            "weight": GSMarenaBaseFormatResult(category="Body", attrs=["Weight"]),
            "jack": GSMarenaBaseFormatResult(category="Sound", attrs=["3.5mm jack"]),
            "usb": GSMarenaBaseFormatResult(category="Comms", attrs=["USB"]),
            "sensors": GSMarenaBaseFormatResult(category="Features", attrs=["Sensors"]),
            "battery": GSMarenaBaseFormatResult(category="Battery", attrs=["Type"]),
            "charging": GSMarenaBaseFormatResult(category="Battery", attrs=["Charging"]),
            "display": GSMarenaBaseFormatResult(category="Display", attrs=["Type", "Size", "Resolution"]),
            "chipset": GSMarenaBaseFormatResult(category="Platform", attrs=["Chipset", "CPU", "GPU"]),
            "main_camera": GSMarenaBaseFormatResult(category="Main Camera", attrs=[]),
            "selfie_camera": GSMarenaBaseFormatResult(category="Selfie camera", attrs=[]),
            "memory": GSMarenaBaseFormatResult(category="Memory", attrs=["Internal"]),
        }

        results = {}

        for key, specs in specifications_map.items():
            results[key] = self.get_specifications(specifications, specs.category, specs.attrs)

        return results