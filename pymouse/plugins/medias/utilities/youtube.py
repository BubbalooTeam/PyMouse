#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.

from hydrogram.types import CallbackQuery
from pymouse import PyMouse, YT_DLP

class YouTubeUtils:
    @staticmethod
    def YouTubeMakeTextWithInfos(
        Vidurl: str,
        Vidtitle: str,
        Vidduration: str,
        Vidviews: str,
        publishedTime: str,
        CreatorofContentUrl: str,
        CreatorofContentName: str,
        i18n: dict,
    ):
        headInfo = "<b><a href={vidURL}>{vidTitle}</a></b>".format(
            vidURL=Vidurl,
            vidTitle=Vidtitle
        )
        out = i18n["youtube-dl"]["formatter"]["publishedTime"].format(
            publishedTime=publishedTime
        )
        out += i18n["youtube-dl"]["formatter"]["durationTime"].format(
            durationTime=Vidduration
        )
        out += i18n["youtube-dl"]["formatter"]["viewsMedia"].format(
            viewsMedia=Vidviews
        )
        CreatorOfContent = "<a href={channelURL}>{channelName}</a>".format(
            channelURL=CreatorofContentUrl,
            channelName=CreatorofContentName
        )
        out += i18n["youtube-dl"]["formatter"]["creatorOfContent"].format(
            creatorOfContent=CreatorOfContent
        )
        return str(headInfo + out)

    @staticmethod
    async def YouTubeDownloadAndSendMedia(
        c: PyMouse, # type: ignore
        cb: CallbackQuery,
        chat_id: int,
        vid: str,
        thumb: str,
        filepath: str,
        media_type: str,
        format_uid: str,
        i18n: dict
    ):
        # Download and GET infos of the video from YouTube
        yt = await YT_DLP().downloader_method(
            c=PyMouse,
            cb=cb,
            vid=vid,
            filepath=filepath,
            media_type=media_type,
            frt=format_uid,
            i18n=i18n
        )
        # Upload Media to Telegram using method (send_media) of the YT_DLP class
        await YT_DLP().send_media(
            c=c,
            cb=cb,
            chat_id=chat_id,
            thumb=thumb,
            media_type=media_type,
            media_class=yt,
            i18n=i18n,
        )

yt_utils = YouTubeUtils()