package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.Credentials.Companion.CREDENTIALS_TABLE_NAME
import java.util.Base64.getEncoder

@Entity(tableName = CREDENTIALS_TABLE_NAME)
data class Credentials(
        @PrimaryKey
        val login: String,
        val password: String,
        val isAuthorized: Boolean,
        val id: Int?,
) {
    fun toAuthHeader(): String {
        val encodedCredentials = getEncoder().encodeToString("$login:$password".toByteArray())
        return "Basic $encodedCredentials"
    }

    companion object {
        const val CREDENTIALS_TABLE_NAME = "credentials"


    }
}