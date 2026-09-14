package app.everyreview.presentation.nav

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import app.everyreview.presentation.product.ProductScreen
import app.everyreview.presentation.review.WriteReviewScreen
import app.everyreview.presentation.scan.ScanScreen

private const val ROUTE_SCAN = "scan"
private const val ROUTE_PRODUCT = "product/{barcode}"
private const val ROUTE_WRITE_REVIEW = "write_review/{productId}"

@Composable
fun EveryReviewNavHost(navController: NavHostController = rememberNavController()) {
    NavHost(navController = navController, startDestination = ROUTE_SCAN) {
        composable(ROUTE_SCAN) {
            ScanScreen(
                onBarcodeScanned = { barcode -> navController.navigate("product/$barcode") },
            )
        }
        composable(
            ROUTE_PRODUCT,
            arguments = listOf(navArgument("barcode") { type = NavType.StringType }),
        ) {
            ProductScreen(
                onWriteReview = { productId -> navController.navigate("write_review/$productId") },
            )
        }
        composable(
            ROUTE_WRITE_REVIEW,
            arguments = listOf(navArgument("productId") { type = NavType.StringType }),
        ) {
            WriteReviewScreen(
                onSubmitted = { navController.popBackStack() },
            )
        }
    }
}
