package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.ChangePasswordAttempt.Companion.CHANGE_PASSWORD_ATTEMPT_TABLE_NAME

@Entity(tableName = CHANGE_PASSWORD_ATTEMPT_TABLE_NAME)
data class ChangePasswordAttempt(
    @PrimaryKey
    val credentials: String,
    val randomToken: String,
    val expiresInSeconds: Long,
    val createdAt: Long,
) {
    companion object {
        const val CHANGE_PASSWORD_ATTEMPT_TABLE_NAME = "change_password_attempt"
    }
}