from asyncio import sleep
from collections import defaultdict
from json import dumps, loads
from typing import Any, Dict, List, Union, Optional
from tqdm import tqdm
from uuid import uuid4

from yt_dlp import YoutubeDL
from yt_dlp.utils import DownloadError, ExtractorError, GeoRestrictedError

from hydrogram.types import CallbackQuery, InlineKeyboardButton, InlineKeyboardMarkup

from pymouse import Decorators, log
from pymouse.utils import http


from typing import Optional

YT = "https://www.youtube.com/"
YT_VID_URL = YT + "watch?v="

def sublists(input_list: List[Any], width: int = 3) -> List[List[Any]]:
    """retuns a single list of multiple sublist of fixed width"""
    return [input_list[x : x + width] for x in range(0, len(input_list), width)]

def humanbytes(size: float) -> str:
    """humanize size"""
    if not size:
        return ""
    power = 1024
    t_n = 0
    power_dict = {0: " ", 1: "Ki", 2: "Mi", 3: "Gi", 4: "Ti"}
    while size > power:
        size /= power
        t_n += 1
    return "{:.2f} {}B".format(size, power_dict[t_n])

class YouTubeDownloader:
    @staticmethod
    @Decorators().aiowrap()
    def extract_info(
        ins: YoutubeDL,
        url: str,
        download: Optional[bool] = True
    ):
        return ins.extract_info(url, download)

    class YouHooker:
        def __init__(
            self,
            cb: CallbackQuery,
            i18n: dict
        ):
            self.cb = cb
            self.i18n = i18n

        async def youhook(self, d: dict):
            if d.get("status") == "downloading":
                eta = d.get("_eta_str")
                speed = d.get("_speed_str")
                tb = d.get("_downloaded_bytes_str")
                tt = d.get("_total_bytes_str")
                # //
                downloading_text = str(self.i18n["youtube-dl"]["downloading"] + "\n\n" + "Tamanho: {tb} / {tt}" + "\n" + "Velocidade: {speed}" + "\n" + "Aguarde: {eta}").format(
                    tb=tb,
                    tt=tt,
                    speed=speed,
                    eta=eta
                )
                await self.cb.edit_message_caption(caption=downloading_text)

class Buttons(InlineKeyboardMarkup):
    def __init__(self, inline_keyboard: List[List["InlineKeyboardButton"]]):
        super().__init__(inline_keyboard)

    def __add__(self, extra: Union[str, int]) -> InlineKeyboardMarkup:
        """Add extra Data to callback_data of every button

        Parameters:
        ----------
            - extra (`Union[str, int]`): Extra data e.g A `key` or `user_id`.

        Raises:
        ------
            `TypeError`

        Returns:
        -------
            `InlineKeyboardMarkup`: Modified markup
        """
        if not isinstance(extra, (str, int)):
            raise TypeError(
                f"unsupported operand `extra` for + : '{type(extra)}' and '{type(self)}'"
            )
        ikb = self.inline_keyboard
        cb_extra = f"-{extra}"
        for row in ikb:
            for btn in row:
                if (
                    (cb_data := btn.callback_data)
                    and cb_data.startswith("yt_")
                    and not cb_data.endswith(cb_extra)
                ):
                    cb_data += cb_extra
                    btn.callback_data = cb_data[:64]  # limit: 1-64 bytes.
        return InlineKeyboardMarkup(ikb)

    def add(self, extra: Union[str, int]) -> InlineKeyboardMarkup:
        """Add extra Data to callback_data of every button

        Parameters:
        ----------
            - extra (`Union[str, int]`): Extra data e.g A `key` or `user_id`.

        Raises:
        ------
            `TypeError`

        Returns:
        -------
            `InlineKeyboardMarkup`: Modified markup
        """
        return self.__add__(extra)


class SearchResult:
    def __init__(
        self,
        key: str,
        text: str,
        image: str,
        buttons: InlineKeyboardMarkup,
    ) -> None:
        self.key = key
        self.buttons = Buttons(buttons.inline_keyboard)
        self.caption = text
        self.image_url = image

    def __repr__(self) -> str:
        out = self.__dict__.copy()
        out["buttons"] = (
            loads(str(btn)) if (btn := out.pop("buttons", None)) else None
        )
        return dumps(out, indent=4)


class YT_DLP:
    async def get_download_button(self, yt_id: str, user_id: int) -> SearchResult:
        buttons = [
            [
                InlineKeyboardButton(
                    "🥇 BEST - 🎥 MP4",
                    callback_data=f"yt_dl|{yt_id}|mp4|{user_id}|v",
                ),
            ]
        ]
        best_audio_btn = [
            [
                InlineKeyboardButton(
                    "🥇 BEST - 📀 320Kbps - MP3",
                    callback_data=f"yt_dl|{yt_id}|mp3|{user_id}|a",
                )
            ]
        ]

        params = {"no-playlist": True, "quiet": True, "logtostderr": False}

        try:
            # //
            vid_data = await YouTubeDownloader().extract_info(
                YoutubeDL(params), f"{YT_VID_URL}{yt_id}", download=False
            )
        except ExtractorError:
            vid_data = None
            buttons += best_audio_btn
        else:
            # ------------------------------------------------ #
            qual_dict = defaultdict(lambda: defaultdict(int))
            qual_list = ("1440p", "1080p", "720p", "480p", "360p", "240p", "144p")
            audio_dict: Dict[int, str] = {}
            # ------------------------------------------------ #
            for video in vid_data["formats"]:
                fr_note = video.get("format_note")
                fr_id = video.get("format_id")
                fr_size = video.get("filesize")
                if video.get("ext") == "mp4":
                    for frmt_ in qual_list:
                        if fr_note in (frmt_, frmt_ + "140"):
                            qual_dict[frmt_][fr_id] = fr_size
                if video.get("acodec") != "none":
                    bitrrate = video.get("abr")
                    if bitrrate == (0 or "None"):
                        pass
                    else:
                        audio_dict[
                            bitrrate
                        ] = f"📀 {bitrrate}Kbps ({humanbytes(fr_size) or 'N/A'})"
            audio_dict = await self.delete_none(audio_dict)
            video_btns: List[InlineKeyboardButton] = []
            for frmt in qual_list:
                frmt_dict = qual_dict[frmt]
                if len(frmt_dict) != 0:
                    frmt_id = sorted(list(frmt_dict))[-1]
                    frmt_size = humanbytes(frmt_dict.get(frmt_id)) or "N/A"
                    video_btns.append(
                        InlineKeyboardButton(
                            f"🎥 {frmt} ({frmt_size})",
                            callback_data=f"yt_dl|{yt_id}|{frmt_id}+140|{user_id}|v",
                        )
                    )
            buttons += sublists(video_btns, width=2)
            buttons += best_audio_btn
            buttons += sublists(
                list(
                    map(
                        lambda x: InlineKeyboardButton(
                            audio_dict[x], callback_data=f"yt_dl|{yt_id}|{x}|{user_id}|a"
                        ),
                        sorted(audio_dict.keys(), reverse=True),
                    )
                ),
                width=2,
            )

        return SearchResult(
            yt_id,
            (
                f"<a href={YT_VID_URL}{yt_id}>{vid_data.get('title')}</a>"
                if vid_data
                else ""
            ),
            vid_data.get("thumbnail")
            if vid_data
            else "https://s.clipartkey.com/mpngs/s/108-1089451_non-copyright-youtube-logo-copyright-free-youtube-logo.png",
            InlineKeyboardMarkup(buttons),
        )

    @Decorators().aiowrap()
    def delete_none(self, _dict):
        """Delete None values recursively from all of the dictionaries, tuples, lists, sets"""
        if isinstance(_dict, dict):
            for key, value in list(_dict.items()):
                if isinstance(value, (list, dict, tuple, set)):
                    _dict[key] = self.delete_none(value)
                elif value is None or key is None:
                    del _dict[key]

        elif isinstance(_dict, (list, set, tuple)):
            _dict = type(_dict)(self.delete_none(item) for item in _dict if item is not None)

        return _dict

    async def downloader(self, url: str, options: Union[str, Any]):
        try:
            down =  await YouTubeDownloader().extract_info(YoutubeDL(options), url, download=True)
            file = down.get("requested_downloads")[0]["filepath"]
            duration = down.get("duration")
            title = down.get("fulltitle")
            return file, duration, title
        except DownloadError:
            log.error("[DownloadError] : Failed to Download Media")
        except GeoRestrictedError:
            log.error(
                "[GeoRestrictedError] : The uploader has not made this video"
                " available in your country"
            )
        except Exception as e:
            log.error("YouTuber: Something Went Wrong: {}".format(e))

    @Decorators().aiowrap()
    def rand_key(self):
        return str(uuid4())[:8]

    async def get_ytthumb(self, videoid: str):
        thumb_quality = [
            "maxresdefault.jpg",  # Best quality
            "hqdefault.jpg",
            "sddefault.jpg",
            "mqdefault.jpg",
            "default.jpg",  # Worst quality
        ]
        thumb_link = "https://i.imgur.com/4LwPLai.png"
        for qualiy in thumb_quality:
            link = f"https://i.ytimg.com/vi/{videoid}/{qualiy}"
            if (await http.get(link)).status_code == 200:
                thumb_link = link
                break
        return thumb_link