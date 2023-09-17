import { Toolbar, AppBar, Link, Typography } from "@mui/material"
import {useTheme} from "@mui/material/styles";

const PrimayAppBar = () => {
    const theme = useTheme();
    return (
        <AppBar sx={{
            backgroundColor: theme.palette.background.default, 
            borderBottom: `1px solid ${theme.palette.divider}`,
            }}
            >
            <Toolbar variant="dense" sx={{
                height: theme.primaryAppBar.height,
                minHeight: theme.primaryAppBar.height,
                }} 
            >
                <Link href="/" underline="none" color="inherit">
                    <Typography 
                        variant="h4" 
                        noWrap 
                        component="div" 
                        sx={{display:{fontWeight: 700, letterSpacing: "-0.5px"}}}
                    >
                        ChatPjt
                    </Typography>
                </Link>
            </Toolbar>
        </AppBar>
    )
}

export default PrimayAppBar;