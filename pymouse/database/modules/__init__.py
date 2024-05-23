import json
from typing import Optional

from pymouse.utils import UtilsDecorator
from ..exceptions import LocalDatabaseJsonError, LocalDatabaseNotFound

class DataBase:
    def __init__(self):
        self.db_file = "pymouse/database/filepaths/database.json"
        self.load_db()

    def load_db(self):
        try:
            with open(self.db_file, 'r') as file:
                self.db = json.load(file)
        except (LocalDatabaseNotFound, LocalDatabaseJsonError):
            self.db = {}

    def save_db(self):
        with open(self.db_file, 'w') as file:
            json.dump(self.db, file, indent=4)

    class GetCollection:
        def __init__(self, collection: str):
            self.parent_db = DataBase()
            self.collection = collection
            if collection not in self.parent_db.db:
                self.parent_db.db[collection] = []

        @UtilsDecorator().aiowrap()
        def find(self):
            return self.parent_db.db[self.collection]
        
        @UtilsDecorator().aiowrap()
        async def find_one(self, filter: Optional[dict] = None):
            # find global collection info
            collection_data = await self.find()
            if not filter:
                print("[database/modules]: No information provided for find especified info...")
                return
            # find keys
            try:
                index = next((i for i, item in enumerate(collection_data) if all(item[k] == v for k, v in filter.items())), None)
            except (IndexError, KeyError):
                index = None
            # find the information base with index
            if index is not None:
                return collection_data[index]
            else:
                return {}

        @UtilsDecorator().aiowrap()
        async def insert_or_update(self, filter: Optional[dict] = None, info: Optional[dict] = None):
            if not info:
                print("[database/modules]: No information provided for update...")
                return
            inup = False
            collection_data = await self.find()
            if filter:
                # Find the index of the item to be updated
                try:
                    index = next((i for i, item in enumerate(collection_data) if all(item[k] == v for k, v in filter.items())), None)
                except (IndexError, KeyError):
                    index = None
                if index is not None:
                    # Update the existing item
                    collection_data[index].update(info)
                    inup = True
                else:
                    # Insert a new item
                    collection_data.append(info)
                    inup = True

            # Save the changes to the database
            self.parent_db.save_db()
            return inup

        @UtilsDecorator().aiowrap()
        async def delete(self, filter: Optional[dict] = None) -> bool:
            """
            Remove a record from the database based on a specified key and its corresponding value.

            Args:
                `filter (dict)`: A dictionary containing the key-value pair for filtering records.

            Returns:
                `bool`: True if a record was deleted, False otherwise.
            """
            deleted = False
            if filter:
                finder = await self.find()
                for record in finder:
                    if all(record.get(key) == value for key, value in filter.items()):
                        finder.remove(record)
                        deleted = True
            else:
                self.parent_db.db.clear()
                deleted = True
            # Save the changes to the database
            self.parent_db.save_db()
            return deleted