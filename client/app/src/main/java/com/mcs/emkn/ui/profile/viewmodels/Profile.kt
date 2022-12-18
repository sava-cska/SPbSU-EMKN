package com.mcs.emkn.ui.profile.viewmodels

import android.net.Uri
import android.os.Parcel
import android.os.Parcelable
import retrofit2.http.Url

data class Profile(
    val id: Int,
    val avatarUri: Uri,
    val firstName: String,
    val secondName: String,
) : Parcelable {
    constructor(parcel: Parcel) : this(
        parcel.readInt(),
        parcel.readParcelable(Uri::class.java.classLoader) ?: Uri.EMPTY,
        parcel.readString() ?: "",
        parcel.readString() ?: ""
    )

    override fun writeToParcel(parcel: Parcel, flags: Int) {
        parcel.writeInt(id)
        parcel.writeParcelable(avatarUri, flags)
        parcel.writeString(firstName)
        parcel.writeString(secondName)
    }

    override fun describeContents(): Int {
        return 0
    }

    companion object CREATOR : Parcelable.Creator<Profile> {
        override fun createFromParcel(parcel: Parcel): Profile {
            return Profile(parcel)
        }

        override fun newArray(size: Int): Array<Profile?> {
            return arrayOfNulls(size)
        }
    }

}