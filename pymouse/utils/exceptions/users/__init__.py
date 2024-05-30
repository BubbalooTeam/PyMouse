from hydrogram.errors import UserIdInvalid, UserNotParticipant, UsernameInvalid, UsernameNotOccupied, PeerIdInvalid

class UsersError(
    UserIdInvalid,
    UserNotParticipant,
    UsernameInvalid,
    UsernameNotOccupied,
    PeerIdInvalid,
    ):
    """Raise if error in handle users"""