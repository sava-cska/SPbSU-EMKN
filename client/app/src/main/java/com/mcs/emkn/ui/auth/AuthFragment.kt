package com.mcs.emkn.ui.auth

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import com.google.android.material.tabs.TabLayoutMediator
import com.mcs.emkn.R
import com.mcs.emkn.auth.AuthComponent
import com.mcs.emkn.databinding.FragmentAuthBinding
import javax.inject.Inject

class AuthFragment : Fragment() {
    @Inject
    lateinit var authComponent: AuthComponent

    private lateinit var binding: FragmentAuthBinding

    private val titlesRes = listOf(R.string.sign_in_button_text, R.string.sign_up_tab_text)

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        binding = FragmentAuthBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setupViewPager()
        setupTabLayout()
    }

    private fun setupViewPager() {
        binding.authViewPager.adapter = activity?.let { AuthPagerAdapter(it) }
        binding.authViewPager.offscreenPageLimit = titlesRes.size
    }

    private fun setupTabLayout() {
        TabLayoutMediator(binding.authTabLayout, binding.authViewPager) { tab, position ->
            tab.text = resources.getString(titlesRes[position])
        }.attach()
    }
}