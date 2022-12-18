package com.mcs.emkn.ui.emailconfirmation

import android.content.Context
import android.util.AttributeSet
import android.view.LayoutInflater
import androidx.constraintlayout.widget.ConstraintLayout
import com.mcs.emkn.databinding.VerificationCodeSlotViewBinding


class VerificationCodeSlotView @JvmOverloads constructor(
    context: Context,
    attrs: AttributeSet? = null,
    defStyleAttr: Int = 0,
    defStyleRes: Int = 0
) : ConstraintLayout(context, attrs, defStyleAttr, defStyleRes) {

    val viewBinding =
        VerificationCodeSlotViewBinding.inflate(LayoutInflater.from(context), this)
}