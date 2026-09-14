package app.everyreview.domain.usecase

import app.everyreview.domain.model.Product
import app.everyreview.domain.model.ProductDetails
import app.everyreview.domain.repository.ProductRepository
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import io.mockk.slot
import kotlinx.coroutines.test.runTest
import org.junit.Test

class SubmitProductDetailsUseCaseTest {

    private val repository = mockk<ProductRepository>()
    private val useCase = SubmitProductDetailsUseCase(repository)

    @Test
    fun `rejects blank name without calling repository`() = runTest {
        val result = useCase("product-1", ProductDetails(name = "   ", brand = "Acme"))

        assertThat(result.isFailure).isTrue()
        coVerify(exactly = 0) { repository.submitDetails(any(), any()) }
    }

    @Test
    fun `rejects non-numeric or negative price`() = runTest {
        assertThat(useCase("product-1", ProductDetails(name = "Thing", price = "cheap")).isFailure).isTrue()
        assertThat(useCase("product-1", ProductDetails(name = "Thing", price = "-3")).isFailure).isTrue()
        coVerify(exactly = 0) { repository.submitDetails(any(), any()) }
    }

    @Test
    fun `trims fields, drops blanks and upper-cases the currency`() = runTest {
        val expected = Product("product-1", "123", "Thing", "Acme", null, false, 0, null)
        val sent = slot<ProductDetails>()
        coEvery { repository.submitDetails("product-1", capture(sent)) } returns Result.success(expected)

        val result = useCase(
            "product-1",
            ProductDetails(
                name = "  Thing ",
                brand = " Acme ",
                manufacturer = "   ",
                description = "",
                price = " 2,49 ",
                currency = "eur",
            ),
        )

        assertThat(result.getOrNull()).isEqualTo(expected)
        assertThat(sent.captured).isEqualTo(
            ProductDetails(name = "Thing", brand = "Acme", manufacturer = null, description = null, price = "2,49", currency = "EUR"),
        )
    }
}
