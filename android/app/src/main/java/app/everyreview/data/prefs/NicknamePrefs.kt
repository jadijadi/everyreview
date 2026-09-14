package app.everyreview.data.prefs

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.first
import javax.inject.Inject
import javax.inject.Singleton

private val Context.dataStore by preferencesDataStore(name = "everyreview_prefs")

@Singleton
class NicknamePrefs @Inject constructor(
    @ApplicationContext private val context: Context,
) {
    private val nicknameKey = stringPreferencesKey("nickname")

    suspend fun getNickname(): String = context.dataStore.data.first()[nicknameKey] ?: ""

    suspend fun setNickname(nickname: String) {
        context.dataStore.edit { it[nicknameKey] = nickname }
    }
}
