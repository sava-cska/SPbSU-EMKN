package com.mcs.emkn.ui.coursehomeworks

import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import com.mcs.emkn.R
import com.mcs.emkn.core.Router
import com.mcs.emkn.core.rv.RecyclerAdapterWithDelegates
import com.mcs.emkn.core.rv.RouterBundleKeys
import com.mcs.emkn.core.rv.VerticalSpaceDecorator
import com.mcs.emkn.databinding.FragmentCourseHomeworksBinding
import com.mcs.emkn.network.dto.response.HomeworkDto
import com.mcs.emkn.ui.coursehomeworks.viewmodels.CourseHomeworksInteractor
import com.mcs.emkn.ui.coursehomeworks.viewmodels.CourseHomeworksViewModel
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import java.text.SimpleDateFormat
import java.time.Instant
import java.util.Date
import javax.inject.Inject

@AndroidEntryPoint
class CourseHomeworksFragment : Fragment(R.layout.fragment_course_homeworks) {
    @Inject
    lateinit var router: Router

    private var _binding: FragmentCourseHomeworksBinding? = null
    private val binding get() = _binding!!
    private val courseHomeworksInteractor: CourseHomeworksInteractor by viewModels<CourseHomeworksViewModel>()
    private val adapter = RecyclerAdapterWithDelegates(
        listOf(
            CourseHomeworksTaskAdapter(),
            CourseHomeworksSectionAdapter()
        ),
        listOf(
        )
    )

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentCourseHomeworksBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        subscribeToLoadHomeworksStatus()
        arguments?.getInt(RouterBundleKeys.courseId)?.let {
            courseHomeworksInteractor.loadHomeworks(it)
        }
        binding.courseName.text = arguments?.getString(RouterBundleKeys.title) ?: "..."
        binding.coursesRecycler.adapter = adapter
        binding.coursesRecycler.addItemDecoration(
            VerticalSpaceDecorator(
                view.context.resources.getDimensionPixelSize(
                    R.dimen.courses_items_offset
                )
            )
        )
        binding.toCourseDescriptionButtonArrow.setOnClickListener {
            router.back()
        }
    }

    private fun subscribeToLoadHomeworksStatus() {
        lifecycleScope.launch {
            courseHomeworksInteractor.homeworksFlow.collect {
                when (it.hasMore) {
                    true -> {
                        binding.progressBar.visibility = View.VISIBLE
                    }
                    false -> {
                        binding.progressBar.visibility = View.GONE
                        adapter.items = listOf(
                            CourseHomeworksSectionItem(-1, "Открытые задания"),
                        )
                        adapter.items += it.data.homeworks.map { hw ->
                            hw.toCourseHomeworksTaskItem()
                        }
                        adapter.notifyDataSetChanged()
                    }
                }
            }
        }
    }

    fun HomeworkDto.toCourseHomeworksTaskItem(): CourseHomeworksTaskItem {
        val dmy = SimpleDateFormat("dd/MM/yyyy")
        val hm = SimpleDateFormat("HH:mm")
        val dateDeadline = Date(deadline)
        return CourseHomeworksTaskItem(
            id, name, dmy.format(dateDeadline), hm.format(dateDeadline),
            when {
                statusNotPassed != null -> CourseHomeworksTaskItem.HomeworkStatus.NotSubmitted(
                    statusNotPassed.text
                )
                statusUnchecked != null -> CourseHomeworksTaskItem.HomeworkStatus.NotChecked(
                    statusUnchecked.text
                )
                else -> CourseHomeworksTaskItem.HomeworkStatus.Checked(
                    statusChecked?.score ?: 0,
                    statusChecked?.totalScore ?: 0
                )
            }
        )
    }
}