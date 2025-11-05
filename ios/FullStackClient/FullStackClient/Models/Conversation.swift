
import Foundation

struct Conversation: Codable, Identifiable {
    let ID: Int
    let Members: [User]
    let Messages: [Message]
}
