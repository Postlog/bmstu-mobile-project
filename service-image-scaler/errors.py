class ResponseError(Exception):
    pass


class UnexpectedResponseCode(ResponseError):
    pass


class NotFound(ResponseError):
    pass
