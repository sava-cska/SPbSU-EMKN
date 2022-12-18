package com.mcs.emkn.network.dto.errorresponse

import com.mcs.emkn.network.dto.error.NetError
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class RevalidateRegistrationCredentialsErrorResponseErrorsDto(
    @Json(name = "invalid_registration_revalidation") val invalidRegistrationRevalidation: NetError?,
)