// Generated from java-escape by ANTLR 4.11.1
import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link markdownParser}.
 */
public interface markdownListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link markdownParser#file_}.
	 * @param ctx the parse tree
	 */
	void enterFile_(markdownParser.File_Context ctx);
	/**
	 * Exit a parse tree produced by {@link markdownParser#file_}.
	 * @param ctx the parse tree
	 */
	void exitFile_(markdownParser.File_Context ctx);
	/**
	 * Enter a parse tree produced by {@link markdownParser#header}.
	 * @param ctx the parse tree
	 */
	void enterHeader(markdownParser.HeaderContext ctx);
	/**
	 * Exit a parse tree produced by {@link markdownParser#header}.
	 * @param ctx the parse tree
	 */
	void exitHeader(markdownParser.HeaderContext ctx);
	/**
	 * Enter a parse tree produced by {@link markdownParser#list}.
	 * @param ctx the parse tree
	 */
	void enterList(markdownParser.ListContext ctx);
	/**
	 * Exit a parse tree produced by {@link markdownParser#list}.
	 * @param ctx the parse tree
	 */
	void exitList(markdownParser.ListContext ctx);
	/**
	 * Enter a parse tree produced by {@link markdownParser#line}.
	 * @param ctx the parse tree
	 */
	void enterLine(markdownParser.LineContext ctx);
	/**
	 * Exit a parse tree produced by {@link markdownParser#line}.
	 * @param ctx the parse tree
	 */
	void exitLine(markdownParser.LineContext ctx);
}