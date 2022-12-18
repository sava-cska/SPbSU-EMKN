package com.mcs.emkn.ui.tasks

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import com.mcs.emkn.core.Router
import com.mcs.emkn.database.Database
import com.mcs.emkn.databinding.FragmentTasksBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import javax.inject.Inject


@AndroidEntryPoint
class TasksFragment : Fragment() {
    private var _binding: FragmentTasksBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var db: Database

    @Inject
    lateinit var router: Router

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentTasksBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.logoutButton.setOnClickListener {
            logout()
        }
    }

    private fun logout() {
        lifecycleScope.launch(Dispatchers.IO) {
            db.accountsDao().deleteCredentials()
            withContext(Dispatchers.Main) {
                router.goToRegistrationNavGraph()
            }
        }
    }
}