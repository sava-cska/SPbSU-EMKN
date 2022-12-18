package com.mcs.emkn.network.dto.errorresponse

import com.mcs.emkn.network.dto.error.NetError
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class RegistrationErrorResponseErrorsDto(
    @Json(name = "illegal_login") val illegalLogin: NetError?,
    @Json(name = "login_is_not_available") val loginIsNotAvailable: NetError?,
    @Json(name = "illegal_password") val illegalPassword: NetError?,
    @Json(name = "email_is_not_available") val emailIsNotAvailable: NetError?,
)

