from typing import Optional

from pymouse import db, usersmodel_db

class AFKDB:
    async def updateAFK(
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

        await users_db.insert_or_update(
            filter={"user_id": user_id},
            info=afk_mapper,
        )

    async def setAFK(
        self,
        user_id: int,
        time: str,
        reason: Optional[str] = None
    ):
        await self.updateAFK(
            user_id=user_id,
            is_afk=True,
            time=time,
            reason=reason,
        )

    async def unsetAFK(
        self,
        user_id: int,
    ):
        await self.updateAFK(
            user_id=user_id,
            is_afk=False,
        )

    async def getAFK(
        self,
        user_id: int,
    ) -> dict:
        afkmap = await usersmodel_db.users_db.getuser_dict(user_id).get("afk", {})
        return afkmap

   
afk_db = AFKDB()