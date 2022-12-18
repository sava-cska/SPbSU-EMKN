package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.ChangePasswordCommit.Companion.CHANGE_PASSWORD_COMMIT_TABLE_NAME

@Entity(tableName = CHANGE_PASSWORD_COMMIT_TABLE_NAME)
class ChangePasswordCommit(
    @PrimaryKey
    val changePasswordToken: String,
) {
    companion object {
        const val CHANGE_PASSWORD_COMMIT_TABLE_NAME = "change_password_commit"
    }
}