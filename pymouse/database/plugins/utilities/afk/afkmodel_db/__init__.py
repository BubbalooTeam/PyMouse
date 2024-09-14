#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from typing import Optional

from pymouse import db, usersmodel_db

class AFKDB:
    def updateAFK(
        self,
        user_id: int,
        is_afk: bool,
        time: Optional[str] = None,
        reason: Optional[str] = None,
    ):
        users_db = db.GetCollection("users")
        afk_mapper = {
            "user_id": user_id,
            "afk": {
                "is_afk": is_afk,
                "time": time,
                "reason": reason,
            }
        }

        afk_mapper["afk"] = {k: v for k, v in afk_mapper["afk"].items() if v is not None}

        users_db.insert_or_update(
            filter={"user_id": user_id},
            info=afk_mapper,
        )

    def setAFK(
        self,
        user_id: int,
        time: str,
        reason: Optional[str] = None
    ):
        self.updateAFK(
            user_id=user_id,
            is_afk=True,
            time=time,
            reason=reason,
        )

    def unsetAFK(
        self,
        user_id: int,
    ):
        self.updateAFK(
            user_id=user_id,
            is_afk=False,
        )

    def getAFK(
        self,
        user_id: int,
    ) -> dict:
        afkmap = usersmodel_db.users_db.getuser_dict(user_id).get("afk", {})
        return afkmap


afk_db = AFKDB()