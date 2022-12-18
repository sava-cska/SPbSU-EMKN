package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.CheckBoxesStateEntity.Companion.CHECK_BOX_TABLE_NAME


@Entity(tableName = CHECK_BOX_TABLE_NAME)
data class CheckBoxesStateEntity(
    @PrimaryKey
    val isExcludingEnroll: Boolean,
    val isExcludingUnenroll: Boolean,
) {
    companion object {
        const val CHECK_BOX_TABLE_NAME = "check_box"
    }
}