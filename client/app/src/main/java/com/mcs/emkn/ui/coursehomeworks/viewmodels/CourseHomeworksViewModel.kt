package com.mcs.emkn.ui.coursehomeworks.viewmodels

import android.util.Log
import androidx.lifecycle.ViewModel
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.core.State
import com.mcs.emkn.database.Database
import com.mcs.emkn.database.entities.HomeworkEntity
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.CoursesEnrollUnenrollRequestDto
import com.mcs.emkn.network.dto.request.CoursesListRequestDto
import com.mcs.emkn.network.dto.response.CourseHomeworks
import com.mcs.emkn.network.dto.response.HomeworkDto
import com.mcs.emkn.network.dto.response.ScoreStatus
import com.mcs.emkn.network.dto.response.TextStatus
import com.mcs.emkn.ui.courses.viewmodels.CoursesInteractor
import com.mcs.emkn.ui.courses.viewmodels.PeriodCourses
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.asCoroutineDispatcher
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.launch
import java.util.concurrent.Executors
import javax.inject.Inject


@HiltViewModel
class CourseHomeworksViewModel @Inject constructor(
    private val api: Api,
    private val db: Database,
) : CourseHomeworksInteractor, ViewModel() {
    private val _homeworksFlow = MutableSharedFlow<State<CourseHomeworks>>()
    override val homeworksFlow: Flow<State<CourseHomeworks>>
        get() = _homeworksFlow

    private val dispatcher = Executors.newSingleThreadExecutor().asCoroutineDispatcher()
    private val scope = CoroutineScope(dispatcher)

    override fun loadHomeworks(courseId: Int) {
        scope.launch {
            _homeworksFlow.emit(
                State(
                    CourseHomeworks(listOf()), true
                )
            )
            val auth = db.accountsDao().getCredentials().first().toAuthHeader()
            when (val response =
                api.getHomeworks(CoursesEnrollUnenrollRequestDto(courseId), auth)) {
                is NetworkResponse.Success -> {
                    db.coursesDao().putHomeworks(
                        response.body.response.homeworks.map {
                            HomeworkEntity(
                                courseId,
                                it.id,
                                it.name,
                                it.deadline,
                                it.statusNotPassed?.text,
                                it.statusUnchecked?.text,
                                it.statusChecked?.score ?: 0,
                                it.statusChecked?.totalScore ?: 0
                            )
                        }
                    )
                }
                else -> Unit
            }
            emitHomeworks(courseId)
        }
    }

    private suspend fun emitHomeworks(courseId: Int) {
        val homeworks = db.coursesDao().getHomeworksByCourse(courseId)
        _homeworksFlow.emit(
            State(
                CourseHomeworks(
                    homeworks.map {
                        HomeworkDto(
                            it.id,
                            it.name,
                            it.deadline,
                            it.statusNotPassed?.let { t -> TextStatus(t) },
                            it.statusUnchecked?.let { t -> TextStatus(t) },
                            ScoreStatus(it.totalScore, it.score)
                        )
                    }), false
            )
        )
    }
}