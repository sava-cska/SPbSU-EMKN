package com.mcs.emkn.database

import androidx.room.Database
import androidx.room.RoomDatabase
import androidx.room.TypeConverters
import com.mcs.emkn.database.utils.ListIntTypeConverter
import com.mcs.emkn.database.entities.*

@Database(
    entities = [
        Credentials::class,
        SignUpAttempt::class,
        ChangePasswordAttempt::class,
        ChangePasswordCommit::class,
        CourseEntity::class,
        PeriodEntity::class,
        ProfileEntity::class,
        CheckBoxesStateEntity::class,
        HomeworkEntity::class
    ],
    version = 1
)
@TypeConverters(
    ListIntTypeConverter::class
)
abstract class Database : RoomDatabase() {
    abstract fun accountsDao(): AccountsDao
    abstract fun coursesDao(): CoursesDao
}